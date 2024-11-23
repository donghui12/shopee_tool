package service

import (
	"time"
	"shopee_tool/pkg/shopee"
	"gorm.io/gorm"
	"fmt"
	"strings"
	"strconv"
	"shopee_tool/pkg/pool"
	"golang.org/x/sync/errgroup"
	"sync"
)


type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) UpdateOrder(cookies string, day int) error {

	// 2. 构建请求参数
	client := shopee.NewClient(
		shopee.WithTimeout(30*time.Second),
		shopee.WithRetry(3, 5*time.Second),
	)

	session := ""
	sessionList := strings.Split(cookies, ";")
	for _, s := range sessionList {
		if strings.Contains(s, "SPC_CNSC_SESSION") {
			session = s
			break
		}
	}

	session += ";"
	fmt.Printf("session: %s\n", session)

	// 3. 获取店铺列表
	merchantShopList, err := client.GetMerchantShopList(session)
	if err != nil {
		return err
	}
	fmt.Printf("merchant shop list: %v\n", len(merchantShopList))

	// 创建 errgroup
	g := new(errgroup.Group)
	var mu sync.Mutex
	errors := make([]error, 0)

	// 4. 提交任务到全局工作池
	for _, shop := range merchantShopList {
        shop := shop // 创建副本
		// 将 shop.Region 转换为小写
		region := strings.ToLower(shop.Region)
		// 将 shopId 转换为 string
		shopId := strconv.FormatInt(shop.ShopID, 10)
        
        g.Go(func() error {
			executeTask := func() error {
				err := processShopOrder(client, session, region, shopId, day)
				if err != nil {
					mu.Lock()
					errors = append(errors, err)
					mu.Unlock()
				}
				return err
			}
			task := pool.Task{Execute: executeTask}

			pool.GlobalWorkerPool.Submit(task)
			return err
        })
    }

	// 5. 等待所有任务完成
	if err := g.Wait(); err != nil {
		return err
	}

	if len(errors) > 0 {
		return fmt.Errorf("some shops failed: %v", errors)
	}
	fmt.Printf("所有店铺更新完成\n")

	return nil
}

// 处理店铺订单
func processShopOrder(client *shopee.Client, session, region, shopId string, day int) error {
	productIdList, err := client.GetProductList(session, shopId, region)
	if err != nil {
		fmt.Printf("获取商品列表失败: %v\n", err)
		return err
	}
	fmt.Printf("店铺 %s 下. 总共获取商品列表: %d\n", shopId, len(productIdList))

	for _, productId := range productIdList {
		err := client.UpdateProductInfo(productId, day, session, shopId, region)
		if err != nil {
			return err
		}
	}
	return nil
}
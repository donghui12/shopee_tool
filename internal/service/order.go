package service

import (
	"time"
	"shopee_tool/pkg/shopee"
	"gorm.io/gorm"
	"fmt"
	"strings"
	"strconv"
	"shopee_tool/pkg/pool"
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

	// 获取 session
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

	// 4. 获取店铺信息
	totalShopInfoList := make([]shopee.UpdateProductInfoReq, 0)
	for _, shop := range merchantShopList {
		shopId := strconv.FormatInt(shop.ShopID, 10)
		region := strings.ToLower(shop.Region)
		shopInfoList, err := client.GetProductList(session, shopId, region)
		if err != nil {
			return err
		}
		fmt.Printf("shop info list: %v\n", len(shopInfoList))
		for _, productId := range shopInfoList {
			totalShopInfoList = append(totalShopInfoList, shopee.UpdateProductInfoReq{
				ProductId: productId,
				DaysToShip: day,
				Cookies: session,
				ShopID: shopId,
				Region: region,
			})
		}
	}

	pool := pool.GetPool()
	wg := sync.WaitGroup{}
	for _, req := range totalShopInfoList {
		wg.Add(1)
		pool.Submit(func() {
			err := client.UpdateProductInfo(req)
			if err != nil {
				fmt.Printf("更新商品信息: %v 失败: %v\n", req.ProductId, err)
			}
			wg.Done()
		})
	}
	wg.Wait()
	fmt.Printf("所有店铺更新完成\n")

	return nil
}
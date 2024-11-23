package service

import (
	"time"
	"shopee_tool/pkg/shopee"
	"gorm.io/gorm"
	"fmt"
	"strings"
	"strconv"
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

	for _, shop := range merchantShopList {
		// 4. 获取商品列表
		fmt.Printf("获取商店店铺列表: %d\n", shop.ShopID)
		// 将 shop.Region 转换为小写
		region := strings.ToLower(shop.Region)
		// 将 shopId 转换为 string
		shopId := strconv.FormatInt(shop.ShopID, 10)
		productIdList, err := client.GetProductList(session, shopId, region)
		if err != nil {
			fmt.Printf("获取商品列表失败: %v\n", err)
			return err
		}
		fmt.Printf("shop: %s . product list: %v\n", shop.ShopID, len(productIdList))

		for _, productId := range productIdList {
			err := client.UpdateProductInfo(productId, day, session, shopId, region)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

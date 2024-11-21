package service

import (
	"time"
	"shopee_tool/pkg/shopee"
	"gorm.io/gorm"
	"fmt"
)


type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) UpdateOrder(cookies string, day int) error {
	fmt.Printf("update order start\n")

	// 2. 构建请求参数
	client := shopee.NewClient(
		shopee.WithTimeout(30*time.Second),
		shopee.WithRetry(3, 5*time.Second),
	)

	// 3. 获取店铺列表
	merchantShopList, err := client.GetMerchantShopList(cookies)
	if err != nil {
		return err
	}
	fmt.Printf("merchant shop list: %v\n", len(merchantShopList))

	for _, shop := range merchantShopList {
		// 4. 获取商品列表
		productIdList, err := client.GetProductList(cookies, string(shop.ShopID), shop.Region)
		if err != nil {
			return err
		}

		for _, productId := range productIdList {
			err := client.UpdateProductInfo(productId, day, cookies, string(shop.ShopID), shop.Region)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

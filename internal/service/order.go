package service

import (
	"time"
	"shopee_tool/pkg/shopee"
	"gorm.io/gorm"
)


type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) UpdateOrder(cookies string) error {

	// 2. 构建请求参数
	client := shopee.NewClient(
		shopee.WithTimeout(30*time.Second),
		shopee.WithRetry(3, 5*time.Second),
	)

	// 3. 获取商品列表
	productList, err := client.GetProductList(cookies)
	if err != nil {
		return err
	}

	// 4. 遍历商品列表，更新库存
	for _, product := range productList.Data.Products {
		err = client.UpdateProductInfo(product.ProductID, cookies)
		if err != nil {
			return err
		}
	}
	return nil
}

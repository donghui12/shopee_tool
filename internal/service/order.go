package service

import (
	"strings"
	"strconv"
	"sync"

	"gorm.io/gorm"
	"go.uber.org/zap"

	"shopee_tool/pkg/pool"
	"shopee_tool/pkg/shopee"
	"shopee_tool/pkg/logger"
)


type OrderService struct {
	db *gorm.DB
}

func NewOrderService(db *gorm.DB) *OrderService {
	return &OrderService{db: db}
}

func (s *OrderService) UpdateOrder(session string, day int) error {

	// 2. 构建请求参数
	client := shopee.GetShopeeClient()

	logger.Info("Retrieved merchant shop list",
		zap.String("session", session),
	)

	// 3. 获取店铺列表
	merchantShopList, err := client.GetMerchantShopList(session)
	if err != nil {
		return err
	}

	logger.Info("获取店铺列表", zap.Int("shop_count", len(merchantShopList)))

	// 4. 获取店铺信息
	totalShopInfoList := make([]shopee.UpdateProductInfoReq, 0)
	for _, shop := range merchantShopList {
		shopId := strconv.FormatInt(shop.ShopID, 10)
		region := strings.ToLower(shop.Region)
		shopInfoList, err := client.GetProductList(session, shopId, region)
		if err != nil {
			return err
		}
		logger.Info("获取店铺商品列表", zap.String("shop_id", shopId), 
			zap.Int("product_count", len(shopInfoList)))
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
	logger.Info("获取商品列表", zap.Int("total_product_count", len(totalShopInfoList)))

	pool := pool.GetPool()
	wg := sync.WaitGroup{}
	for _, req := range totalShopInfoList {
		wg.Add(1)
		pool.Submit(func() {
			err := client.UpdateProductInfo(req)
			if err != nil {
				logger.Error("更新商品信息失败", zap.Int64("product_id", req.ProductId),
					zap.String("shop_id", req.ShopID), zap.String("region", req.Region),
					zap.Error(err))
			}
			wg.Done()
		})
	}
	wg.Wait()
	logger.Info("所有店铺更新完成")

	return nil
}
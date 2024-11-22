package shopee

// API 路径常量
const (
    // API 基础路径
    BaseSellerURL = "https://seller.shopee.cn"

    // 账户相关接口
    APIPathLogin = "/api/cnsc/selleraccount/login/"

    // 商品相关接口
	APIPathUpdateProductInfo = "/api/v3/product/update_product_info"
    APIPathProductList = "/api/v3/mpsku/list/v2/get_product_list"
    APIPathProductUpdate = "/api/v3/product/update_product"
	APIPathGetMerchantShopList = "/api/cnsc/selleraccount/get_merchant_shop_list/"
	APIPathGetOrSetShop = "/api/cnsc/selleraccount/get_or_set_shop/"
)

// API 请求方法
const (
    HTTPMethodGet  = "GET"
    HTTPMethodPost = "POST"
    HTTPMethodPut  = "PUT"
)

// API 响应码
const (
    ResponseCodeSuccess = 0
    ResponseCodeError   = 1
) 
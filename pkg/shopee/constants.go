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
    APIPathProductDetail = "/api/v3/product/get_product_detail"
    APIPathProductUpdate = "/api/v3/product/update_product"
    
    // 订单相关接口
    APIPathOrderList = "/api/v3/order/get_order_list"
    APIPathOrderDetail = "/api/v3/order/get_order_detail"
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
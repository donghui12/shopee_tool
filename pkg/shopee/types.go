package shopee

// LoginRequest 登录请求参数
type LoginRequest struct {
    PasswordHash    string `json:"password_hash"`
    Remember        bool   `json:"remember"`
    OtpType        int    `json:"otp_type"`
    SubaccountPhone string `json:"subaccount_phone"`
    Vcode          string `json:"vcode,omitempty"`
}

// LoginResponse 登录响应
type LoginResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Token string `json:"token"`
    } `json:"data"`
}

// ProductListRequest 商品列表请求
type ProductListRequest struct {
    PageSize   int    `json:"page_size"`
    PageNo     int    `json:"page_no"`
    SearchType string `json:"search_type,omitempty"`
    Keyword    string `json:"keyword,omitempty"`
    SortBy     string `json:"sort_by,omitempty"`
    SortType   int    `json:"sort_type,omitempty"`
}

type OngoingCampaigns struct {
    ProductID    int `json:"product_id"`
}

type PromotionDetail struct {
    OngoingCampaigns []OngoingCampaigns `json:"ongoing_campaigns"`
}

type Product struct {
    PromotionDetail PromotionDetail `json:"promotion_detail"`
}

type PageInfo struct {
    PageNumber int `json:"page_number"`
    PageSize   int `json:"page_size"`
    Total      int `json:"total"`
}

type ProductListData struct {
    Products []Product `json:"products"`
    PageInfo PageInfo  `json:"page_info"`
}

// ProductListResponse 商品列表响应
type ProductListResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
	UserMessage string `json:"user_message"`
    Data    ProductListData `json:"data"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
    ProductID   int64   `json:"product_id"`
    Name        string  `json:"name,omitempty"`
    Description string  `json:"description,omitempty"`
    Price       float64 `json:"price,omitempty"`
    Stock       int     `json:"stock,omitempty"`
}

type MerchantShop struct {
    Region string `json:"region"`
    ShopID int64 `json:"shop_id"`
}

type MerchantShopList struct {
	Shops []MerchantShop `json:"shops"`
}

type MerchantShopListResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    MerchantShopList `json:"data"`
}
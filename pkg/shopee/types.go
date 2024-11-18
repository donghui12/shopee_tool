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

// ProductListResponse 商品列表响应
type ProductListResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Data    struct {
        Products []Product `json:"products"`
        Total    int      `json:"total"`
    } `json:"data"`
}

// Product 商品信息
type Product struct {
    ProductID   int64   `json:"product_id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
    CreateTime  int64   `json:"create_time"`
    UpdateTime  int64   `json:"update_time"`
}

// UpdateProductRequest 更新商品请求
type UpdateProductRequest struct {
    ProductID   int64   `json:"product_id"`
    Name        string  `json:"name,omitempty"`
    Description string  `json:"description,omitempty"`
    Price       float64 `json:"price,omitempty"`
    Stock       int     `json:"stock,omitempty"`
}
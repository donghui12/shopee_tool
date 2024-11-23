package shopee

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
	"net/url"
	"strconv"
	"github.com/google/uuid"
)

type Client struct {
    baseURL     string
    httpClient  *http.Client
    cookies     []*http.Cookie
    userAgent   string
    retryTimes  int
    retryDelay  time.Duration
}

type ClientOption func(*Client)

// WithBaseURL 设置基础URL
func WithBaseURL(url string) ClientOption {
    return func(c *Client) {
        c.baseURL = url
    }
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
    return func(c *Client) {
        c.httpClient.Timeout = timeout
    }
}

// WithRetry 设置重试次数和延迟
func WithRetry(times int, delay time.Duration) ClientOption {
    return func(c *Client) {
        c.retryTimes = times
        c.retryDelay = delay
    }
}

// NewClient 创建新的客户端
func NewClient(options ...ClientOption) *Client {
    client := &Client{
        baseURL: BaseSellerURL,
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        userAgent:   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
        retryTimes:  3,
        retryDelay:  5 * time.Second,
    }

    for _, option := range options {
        option(client)
    }

    return client
}

// Login 登录
func (c *Client) Login(phone, password, vcode string) (string, error) {
	cookieString := "" 
    // 构建表单数据
	urlValue := url.Values{}
	urlValue.Set("password_hash", MD5Hash(password))
	urlValue.Set("remember", "false")
	urlValue.Set("otp_type", "-1")
	urlValue.Set("subaccount_phone", formatPhone(phone))
	if vcode != "" {
		urlValue.Set("vcode", vcode)
		urlValue.Set("otp_type", "-1")
	}
	reqBody := urlValue.Encode()

    // 创建请求
    req, err := http.NewRequest(HTTPMethodPost, c.baseURL+APIPathLogin, strings.NewReader(reqBody))
    if err != nil {
        return cookieString, fmt.Errorf("create login request failed: %w", err)
    }

    // 设置表单请求头
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

    // 执行请求
    resp, err := c.executeWithRetry(req)
    if err != nil {
        return cookieString, fmt.Errorf("login request failed: %w", err)
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return cookieString, fmt.Errorf("read login response failed: %w", err)
    }

    // 解析响应
    var loginResp LoginResponse
    if err := json.Unmarshal(body, &loginResp); err != nil {
        return cookieString, fmt.Errorf("parse login response failed: %w", err)
    }

    // 检查响应状态
    if loginResp.Code != ResponseCodeSuccess {
        return cookieString, fmt.Errorf("login failed: %s", loginResp.Message)
    }
	if loginResp.Message == "error_need_vcode"{
		return cookieString, fmt.Errorf("需要验证码")
	}
	if loginResp.Message == "error_invalid_vcode" {
		return cookieString, fmt.Errorf("验证码错误")
	}
	if loginResp.Message == "error_name_or_password_incorrect" {
		return cookieString, fmt.Errorf("账号或密码错误")
	}

	cookies := resp.Header["Set-Cookie"]
	// 将 cookie 转换为字符串
	for _, cookie := range cookies {
		if strings.Contains(cookie, "SPC_CDS") {
			continue
		}
		cookie = strings.Split(cookie, ";")[0]
		cookieString += cookie + "; "
	}

    return cookieString, nil
}

// GetCookies 获取cookies
func (c *Client) GetCookies() string {
	// 转换cookie string
	cookieJSON, err := json.Marshal(c.cookies)
	if err != nil {
		return ""
	}
	return string(cookieJSON)
}

// get_merchant_shop_list
func (c *Client) GetMerchantShopList(cookies string) ([]MerchantShop, error) {
	merchantShopListResp := &MerchantShopListResponse{}
	resp, err := c.doRequest(HTTPMethodGet, APIPathGetMerchantShopList, nil, cookies)
	if err != nil {
		return nil, fmt.Errorf("get merchant shop list failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body failed: %w", err)
	}
	err = json.Unmarshal(body, &merchantShopListResp)
	if err != nil {
		return nil, fmt.Errorf("unmarshal merchant shop list response failed: %w", err)
	}
	return merchantShopListResp.Data.Shops, nil
}

// GetProductList 获取商品列表
func (c *Client) GetProductList(cookies, shopID, region string) ([]int64, error) {
	fmt.Printf("cookies: %s\n", cookies)
	fmt.Printf("shopID: %s\n", shopID)
	fmt.Printf("region: %s\n", region)

	var productIDs []int64
	var productIDMap = make(map[int64]bool)

	SPC_CDS := uuid.New().String()
	cookies += "SPC_CDS=" + SPC_CDS + ";"
	getProductListParams := url.Values{}
	getProductListParams.Set("SPC_CDS", SPC_CDS)
	getProductListParams.Set("SPC_CDS_VER", "2")
	getProductListParams.Set("list_type", "all")
	getProductListParams.Set("need_ads", "true")
	getProductListParams.Set("cnsc_shop_id", shopID)
	getProductListParams.Set("cbsc_shop_region", region)

	pageNumber := 1
	pageSize := 48
	total := 0
	for {
		productListResp := &ProductListResponse{}
		getProductListParams.Set("page_number", strconv.Itoa(pageNumber))
		getProductListParams.Set("page_size", strconv.Itoa(pageSize))

		APIProductList := APIPathProductList + "?" + getProductListParams.Encode()

		resp, err := c.doRequest(HTTPMethodGet, APIProductList, nil, cookies)
		if err != nil {
			return nil, fmt.Errorf("get product list failed: %w", err)
		}

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response body failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("get product list failed: %s", string(body))
		}

		err = json.Unmarshal(body, &productListResp)
		if err != nil {
			return nil, fmt.Errorf("unmarshal product list response failed: %w", err)
		}

		for _, product := range productListResp.Data.Products {
			for _, campaign := range product.PromotionDetail.OngoingCampaigns {
				productIDMap[int64(campaign.ProductID)] = true
			}
		}

		total += len(productListResp.Data.Products)
		if total >= productListResp.Data.PageInfo.Total {
			break
		}
		pageNumber++
	}

	for productID := range productIDMap {
		productIDs = append(productIDs, productID)
	}

	return productIDs, nil
}

// UpdateProductInfoRequest 更新商品信息请求
type UpdateProductInfoRequest struct {
    ProductID int64       `json:"product_id"`
    ProductInfo ProductInfo `json:"product_info"`
    IsDraft   bool       `json:"is_draft"`
}

// ProductInfo 商品信息
type ProductInfo struct {
    EnableModelLevelDts bool         `json:"enable_model_level_dts"`
    PreOrderInfo       PreOrderInfo `json:"pre_order_info"`
}

// PreOrderInfo 预售信息
type PreOrderInfo struct {
    PreOrder    bool `json:"pre_order"`
    DaysToShip  int  `json:"days_to_ship"`
}

// UpdateProductInfoResponse 更新商品信息响应
type UpdateProductInfoResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	UserMessage string `json:"user_message"`
	Data struct {
		ProductID int64 `json:"product_id"`
	} `json:"data"`
}

// UpdateProductInfo 更新商品信息
func (c *Client) UpdateProductInfo(productID int64, day int, cookies, shopID,region string) error {
	SPC_CDS := uuid.New().String()
	cookies += "SPC_CDS=" + SPC_CDS + ";"

	updateProductInfoParams := url.Values{}
	updateProductInfoParams.Set("SPC_CDS", SPC_CDS)
	updateProductInfoParams.Set("SPC_CDS_VER", "2")
	updateProductInfoParams.Set("cnsc_shop_id", shopID)
	updateProductInfoParams.Set("cbsc_shop_region", region)

	req := &UpdateProductInfoRequest{
		ProductID: productID,
		ProductInfo: ProductInfo{
			EnableModelLevelDts: false,
			PreOrderInfo: PreOrderInfo{
				PreOrder:    true,
				DaysToShip:  day,
			},
		},
		IsDraft:   false,
	}

	APIUpdateProductInfo := APIPathUpdateProductInfo + "?" + updateProductInfoParams.Encode()
	resp, err := c.doRequest(HTTPMethodPost, APIUpdateProductInfo, req, cookies)
	if err != nil {
		return fmt.Errorf("update product info failed: %w", err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("update product info failed: %s", string(body))
	}
	var updateProductInfoResp UpdateProductInfoResponse
	err = json.Unmarshal(body, &updateProductInfoResp)
	if err != nil {
		return fmt.Errorf("unmarshal update product info response failed: %w", err)
	}
	if updateProductInfoResp.Code != ResponseCodeSuccess {
		return fmt.Errorf("update product info failed: %s", updateProductInfoResp.Message)
	}
	fmt.Printf("update product info success: %v\n", updateProductInfoResp.Data.ProductID)
	
	return nil
}

// GetOrSetShopResponse 获取或设置店铺响应
type GetOrSetShopResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data struct {
		ShopID string `json:"shop_id"`
	} `json:"data"`
}

// curl 'https://seller.shopee.cn/api/cnsc/selleraccount/get_or_set_shop/?SPC_CDS=dbb46082-6eb0-473e-b6b9-42fe3be4103a&SPC_CDS_VER=2' \
//   -H 'cookie: SPC_CDS=dbb46082-6eb0-473e-b6b9-42fe3be4103a; SPC_F=nESjUPBIRRyYopsGoJBhvCXpmStUnvLE; _QPWSDCXHZQA=73431902-f460-4a8c-8a7b-0227ccea192c; REC7iLP4Q=f2d09793-28b5-472a-b3bb-a9b10e49755a; SPC_CDS_CHAT=24f6c73f-f519-4680-8f0b-0297bd9d9d57; _sapid=6929d8fe0c56e7996c7414fcaf5271b5d13f9d1feae8a09bfad825fa; shopee_webUnique_ccd=8Etf2zGHigzjyP0AlZoaMQ%3D%3D%7CBs5vGnFM1HVL%2FHneBxdqub4pNjiMjkSje54zMErTYi9Dzbz6Bj5BgCfRIZzhjFUpcJKQFhQ3w36%2FnyQ0us0%3D%7CAcZgRkkzX%2FlXcy43%7C08%7C3; ds=5ed74830f63924150198cd2741f36563; UYOMAPJWEMDGJ=; SPC_CNSC_SESSION=670a7546a7c642bfb1706cc7bf72e946_2_2375038; SPC_SC_OFFLINE_TOKEN=eyJkYXRhIjoib2o4R3ZZN1djOGt4djhOL2x3djVqcjhzN2VsaXprVTBrVFY0eG5JeDlnQUhSWGt6ZDV2RmhCTDYyaGZ4MWU3enR0NWxidksySDVPWmY3MjZaajJxakhId2c4c3FrUVNOd2c1cWNFa0toSFFWQVo5anBaZ0xMbmo0c0kvcEJtdlZTOVZLSnBJTkNHbWNkazdaa0Q5S2pnPT0iLCJpdiI6IkRFaXlYcTI3TkxNRlN4eXhJcXIyaXc9PSIsInNpZ24iOiI1cVB2S2xiTmlSN0FnclNSSFdYdDRJRm9yT0xqUG5FdjlLOUpRbTVTSm5BbGJRdmw3TU42Q3I3VTliVlIrK0h2b3lvb1ZIaTZtVCtpL1kxV3JwOUhBUT09In0=; SPC_CNSC_TK=d4a06e168e89374680af0c0bbe7940cc; SPC_CNSC_UD=1_2375038; CNSC_SSO=TnB2WkNxNjdZS0VIanQ5SpRPARcuQKiFyEEdsBDcXMQJwwAnSOpdNgHSiObLuQFR' \
//   --data-raw '{}' \
//   --compressed
func (c *Client) GetOrSetShop(cookies string) error {
	// 构建请求
	SPC_CDS := uuid.New().String()
	cookies += "SPC_CDS=" + SPC_CDS + ";"

	APIGetOrSetShop := APIPathGetOrSetShop + "?" + url.Values{"SPC_CDS": {SPC_CDS}, "SPC_CDS_VER": {"2"}}.Encode()
	respBody := map[string]interface{}{}
	resp, err := c.doRequest(HTTPMethodPost, APIGetOrSetShop, respBody, cookies)
	if err != nil {
		fmt.Printf("get or set shop request failed: %s\n", err)
		return fmt.Errorf("get or set shop request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read response body failed: %s\n", err)
		return fmt.Errorf("read response body failed: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("get or set shop request failed: %s\n", string(body))
		return fmt.Errorf("get or set shop request failed: %s", string(body))
	}
	var getOrSetShopResp GetOrSetShopResponse
	err = json.Unmarshal(body, &getOrSetShopResp)
	if err != nil {
		fmt.Printf("unmarshal get or set shop response failed: %s\n", err)
		return fmt.Errorf("unmarshal get or set shop response failed: %w", err)
	}
	if getOrSetShopResp.Message != "success" {
		fmt.Printf("get or set shop response failed: %s\n", getOrSetShopResp.Message)
		return fmt.Errorf("get or set shop response failed: %s", getOrSetShopResp.Message)	
	}
	return nil

}

// doRequest 执行请求
func (c *Client) doRequest(method, path string, reqBody interface{}, cookies string) (*http.Response, error) {
    url := c.baseURL + path

    var bodyReader io.Reader
    if reqBody != nil {
        jsonData, err := json.Marshal(reqBody)
        if err != nil {
            return nil, fmt.Errorf("marshal request body failed: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("create request failed: %w", err)
    }

    // 设置 JSON 请求头
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookies)
	req.Header.Set("Host", "seller.shopee.cn")
    c.setCommonHeaders(req)

    resp, err := c.executeWithRetry(req)
    if err != nil {
        return nil, fmt.Errorf("execute request failed: %w", err)
    }

    return resp, nil
}

// handleResponse 处理响应
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("read response body failed: %w", err)
    }

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
    }

    if result == nil {
        return nil
    }

	var productListResp ProductListResponse

    if err := json.Unmarshal(body, &productListResp); err != nil {
        return fmt.Errorf("unmarshal response failed: %w", err)
    }
	fmt.Printf("productListResp: %v\n", productListResp.Data.Products)

    return nil
}
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
func (c *Client) Login(phone, password, vcode string) error {
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
        return fmt.Errorf("create login request failed: %w", err)
    }

    // 设置表单请求头
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

    // 执行请求
    resp, err := c.executeWithRetry(req)
    if err != nil {
        return fmt.Errorf("login request failed: %w", err)
    }
    defer resp.Body.Close()

    // 读取响应
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("read login response failed: %w", err)
    }

    // 解析响应
    var loginResp LoginResponse
    if err := json.Unmarshal(body, &loginResp); err != nil {
        return fmt.Errorf("parse login response failed: %w", err)
    }

    // 检查响应状态
    if loginResp.Code != ResponseCodeSuccess {
        return fmt.Errorf("login failed: %s", loginResp.Message)
    }
	if loginResp.Message == "error_need_vcode" || loginResp.Message == "error_invalid_vcode" {
		return fmt.Errorf("login failed: %s", loginResp.Message)
	}

    // 保存 cookies
    c.cookies = resp.Cookies()

    return nil
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

// GetProductList 获取商品列表
func (c *Client) GetProductList(cookies string) (*ProductListResponse, error) {
    var resp ProductListResponse
    err := c.doRequest(HTTPMethodPost, APIPathProductList, nil, &resp, cookies)
    if err != nil {
        return nil, fmt.Errorf("get product list failed: %w", err)
    }

    if resp.Code != ResponseCodeSuccess {
        return nil, fmt.Errorf("get product list failed: %s", resp.Message)
    }

    return &resp, nil
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
func (c *Client) UpdateProductInfo(productID int64, cookies string) error {
	req := &UpdateProductInfoRequest{
		ProductID: productID,
		ProductInfo: ProductInfo{
			EnableModelLevelDts: false,
			PreOrderInfo: PreOrderInfo{
				PreOrder:    true,
				DaysToShip:  10,
			},
		},
		IsDraft:   false,
	}

	var resp UpdateProductInfoResponse
	err := c.doRequest(HTTPMethodPost, APIPathUpdateProductInfo, req, &resp, cookies)
	if err != nil {
		return fmt.Errorf("update product info failed: %w", err)
	}
	return nil
}

// doRequest 执行请求
func (c *Client) doRequest(method, path string, reqBody interface{}, respBody interface{}, cookies string) error {
    url := c.baseURL + path

    var bodyReader *bytes.Buffer
    if reqBody != nil {
        jsonData, err := json.Marshal(reqBody)
        if err != nil {
            return fmt.Errorf("marshal request body failed: %w", err)
        }
        bodyReader = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return fmt.Errorf("create request failed: %w", err)
    }

    // 设置 JSON 请求头
    req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", cookies)
    c.setCommonHeaders(req)

    resp, err := c.executeWithRetry(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    return c.handleResponse(resp, respBody)
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

    if err := json.Unmarshal(body, result); err != nil {
        return fmt.Errorf("unmarshal response failed: %w", err)
    }

    return nil
}
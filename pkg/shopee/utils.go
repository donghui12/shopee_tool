package shopee

import (
	"crypto/md5"
	"fmt"
	"strings"
	"net/http"
	"time"
)

// MD5Hash 计算MD5哈希值
func MD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return fmt.Sprintf("%x", hash)
}

// formatPhone 格式化手机号
func formatPhone(phone string) string {
	phone = strings.TrimSpace(phone)
	if !strings.HasPrefix(phone, "86") {
		return "86" + phone
	}
	return phone
}

// setCommonHeaders 设置通用请求头
func (c *Client) setCommonHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Origin", BaseSellerURL)
	req.Header.Set("Referer", BaseSellerURL)

	// 添加 cookies
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}
}

// executeWithRetry 带重试的请求执行
func (c *Client) executeWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.retryTimes; i++ {
		resp, err = c.httpClient.Do(req)
		if err == nil {
			return resp, nil
		}
		if i == c.retryTimes {
			return resp, fmt.Errorf("request failed after %d retries: %w", c.retryTimes, err)
		}
		time.Sleep(c.retryDelay)
	}

	return resp, err
}
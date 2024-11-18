package auth

import (
	"shopee_tool/internal/database/models"
	"shopee_tool/pkg/shopee"
	"time"
	"fmt"
	"gorm.io/gorm"
	"encoding/json"
)

type LoginService struct {
	db *gorm.DB
}

func NewLoginService(db *gorm.DB) *LoginService {
	return &LoginService{db: db}
}

type CookieInfo struct {
	Name     string `json:"name"`
	Value    string `json:"value"`
	Domain   string `json:"domain"`
	Path     string `json:"path"`
	Expires  string `json:"expires"`
	HttpOnly bool   `json:"http_only"`
	Secure   bool   `json:"secure"`
}

func (s *LoginService) Login(username, password, vcode string) ([]CookieInfo, error) {
	client := shopee.NewClient(
		shopee.WithTimeout(30*time.Second),
		shopee.WithRetry(3, 5*time.Second),
	)
	
	// 执行登录
	err := client.Login(username, password, vcode)
	if err != nil {
		return nil, err
	}
	
	// 获取cookies
	cookies := client.GetCookies()
	
	// 转换cookie格式
	var cookieInfos []CookieInfo
	for _, cookie := range cookies {
		cookieInfo := CookieInfo{
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			Expires:  cookie.Expires.Format(time.RFC3339),
			HttpOnly: cookie.HttpOnly,
			Secure:   cookie.Secure,
		}
		cookieInfos = append(cookieInfos, cookieInfo)
	}
	
	// 保存到数据库
	if err := s.saveCookies(username, cookieInfos); err != nil {
		return nil, err
	}
	
	return cookieInfos, nil
}

func (s *LoginService) saveCookies(username string, cookies []CookieInfo) error {
	// TODO: 实现保存cookie到数据库的逻辑
	// 将cookie转换为JSON字符串
	cookieJSON, err := json.Marshal(cookies)
	if err != nil {
		return fmt.Errorf("序列化cookie失败: %w", err)
	}

	// 更新数据库中的cookie记录
	result := s.db.Model(&models.Cookie{}).
		Where("username = ?", username).
		Update("cookies", string(cookieJSON))
	
	if result.Error != nil {
		return fmt.Errorf("保存cookie失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("未找到用户: %s", username)
	}
	return nil
}
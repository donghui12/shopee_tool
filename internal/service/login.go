package service

import (
	"shopee_tool/internal/database/models"
	"shopee_tool/pkg/shopee"
	"time"
	"fmt"
	"gorm.io/gorm"
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

func (s *LoginService) Login(username, password, vcode string) error {
	client := shopee.NewClient(
		shopee.WithTimeout(30*time.Second),
		shopee.WithRetry(3, 5*time.Second),
	)
	
	// 执行登录
	cookies, err := client.Login(username, password, vcode)
	if err != nil {
		return err
	}

	err = client.GetOrSetShop(cookies)
	if err != nil {
		return err
	}

	// 创建账户
	err = s.createAccount(username, cookies)
	if err != nil {
		return err
	}
	
	return nil
}

func (s *LoginService) createAccount(username string, cookies string) error {
	// 创建账户
	account := models.Account{
		Username: username,
		Cookies:  cookies,
	}

	// 保存到数据库, 如果账户已存在则更新
	var existingAccount models.Account
	result := s.db.Where("username = ?", username).First(&existingAccount)
	if result.Error == nil {
		// 更新账户
		existingAccount.Cookies = cookies
		s.db.Save(&existingAccount)
		return nil
	}

	// 创建账户
	result = s.db.Create(&account)
	if result.Error != nil {
		return fmt.Errorf("创建账户失败: %w", result.Error)
	}

	return nil
}
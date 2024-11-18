package auth

import (
	"net/http"
	"shopee_tool/internal/database/models"
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

func (s *LoginService) Login(username, password string) ([]CookieInfo, error) {
	// TODO: 实现实际的虾皮登录逻辑
	// 这里是示例代码
	cookies := []CookieInfo{
		{
			Name:     "SPC_EC",
			Value:    "example_cookie_value",
			Domain:   ".shopee.com",
			Path:     "/",
			Expires:  "2024-12-31T23:59:59Z",
			HttpOnly: true,
			Secure:   true,
		},
	}

	// 保存账号信息
	account := &models.Account{
		Username: username,
		Password: password, // 注意：实际使用时需要加密存储
	}
	
	if err := s.db.Create(account).Error; err != nil {
		return nil, err
	}

	// 保存 cookie 信息
	for _, cookie := range cookies {
		dbCookie := &models.Cookie{
			AccountID: account.ID,
			Name:     cookie.Name,
			Value:    cookie.Value,
			Domain:   cookie.Domain,
			Path:     cookie.Path,
			// 需要将 cookie.Expires 字符串解析为 time.Time
		}
		
		if err := s.db.Create(dbCookie).Error; err != nil {
			return nil, err
		}
	}

	return cookies, nil
} 
package service

import (
	"shopee_tool/internal/database/models"
	"time"
	"errors"
	"gorm.io/gorm"
)

type ActiveCodeService struct {
	db *gorm.DB
}

func NewActiveCodeService(db *gorm.DB) *ActiveCodeService {
	return &ActiveCodeService{db: db}
}

func (s *ActiveCodeService) CreateActiveCode(code, expiredAt string) (string, error) {
	if code == "" || expiredAt == "" {
		return "", errors.New("code and expiredAt cannot be empty")
	}
	// 将 expiredAt 转换为 time.Time
	expiredAtTime, err := time.Parse(time.RFC3339, expiredAt)
	if err != nil {
		return "", err
	}
	activeCode := models.ActiveCode{
		Code: code,
		ExpiredAt: expiredAtTime,
	}
	result := s.db.Create(&activeCode)
	if result.Error != nil {
		return "", result.Error
	}
	return activeCode.Code, nil
}

func (s *ActiveCodeService) GetActiveCode(code string) (string, error) {
	var activeCode models.ActiveCode
	result := s.db.Where("code = ?", code).First(&activeCode)
	if result.Error != nil {
		return "", result.Error
	}
	return activeCode.Code, nil
}
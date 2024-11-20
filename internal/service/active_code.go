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

	expiredAtTime, err := time.Parse("2006-01-02 15:04:05", expiredAt)
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

func (s *ActiveCodeService) GetActiveCode(code string) (int, error) {
	var activeCode models.ActiveCode
	result := s.db.Where("code = ? and expired_at > ?", code, time.Now()).First(&activeCode)
	if result.Error != nil {
		return 0, result.Error
	}
	// 查看距离过期还有多少天
	day := activeCode.ExpiredAt.Sub(time.Now()).Hours() / 24
	return int(day), nil
}
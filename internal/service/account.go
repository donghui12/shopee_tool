package service

import (
	"shopee_tool/internal/database/models"
	"gorm.io/gorm"
)

type AccountService struct {
	db *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{db: db}
}

func (s *AccountService) UpdateMachineCode(username, MachineCode string) error {
	// 更新数据库中的机器码
	result := s.db.Model(&models.Account{}).
		Where("username = ?", username).
		Update("machine_code", MachineCode)
	return result.Error
}

func (s *AccountService) GetMachineCode(username string) (string, error) {
	var machineCode string
	result := s.db.Model(&models.Account{}).
		Where("username = ?", username).
		Select("machine_code").Scan(&machineCode)
	return machineCode, result.Error
}

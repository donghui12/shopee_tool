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

func (s *AccountService) GetMachineCode(username, machineCode string) error {
	var selectMachineCode string
	result := s.db.Model(&models.Account{}).
		Where("username = ? AND machine_code = ?", username, machineCode).
		Select("machine_code").Scan(&selectMachineCode)
	return result.Error
}

func (s *AccountService) GetActiveCodeByActiveCode(activeCode string) error {
	var selectActiveCode string
	result := s.db.Model(&models.Account{}).
		Where("active_code = ?", activeCode).
		Select("active_code").Scan(&selectActiveCode)
	return result.Error
}

func (s *AccountService) UpdateActiveCode(username, activeCode string) error {
	result := s.db.Model(&models.Account{}).
		Where("username = ?", username).
		Update("active_code", activeCode)
	return result.Error
}

func (s *AccountService) GetActiveCode(username string) (string, error) {
	var activeCode string
	result := s.db.Model(&models.Account{}).
		Where("username = ?", username).
		Select("active_code").Scan(&activeCode)
	return activeCode, result.Error
}

func (s *AccountService) GetCookies(username string) (string, error) {
    var account models.Account
    result := s.db.Where("username = ?", username).First(&account)
    if result.Error != nil {
        return "", result.Error
    }
    return account.Cookies, nil
}

func (s *AccountService) GetSession(username string) (string, error) {
    var account models.Account
    result := s.db.Where("username = ?", username).First(&account)
    if result.Error != nil {
        return "", result.Error
    }
    return account.Session, nil
}

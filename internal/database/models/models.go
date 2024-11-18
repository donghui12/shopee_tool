package models

import "time"

// Account 存储虾皮账号信息
type Account struct {
    ID        uint      `gorm:"primarykey"`
    Username  string    `gorm:"type:varchar(255);uniqueIndex;not null"` // 修改为 varchar 类型
    Password  string    `gorm:"type:varchar(255);not null"`             // 修改为 varchar 类型
	Phone     string    `gorm:"type:varchar(255);not null"`             // 修改为 varchar 类型
	MachineCode string    `gorm:"type:varchar(255);not null"`             // 修改为 varchar 类型
    CreatedAt time.Time `gorm:"index"`
    UpdatedAt time.Time
}

// Cookie 存储登录后的 cookie 信息
type Cookie struct {
    ID        uint      `gorm:"primaryKey"`
    AccountID uint      `gorm:"index;not null"`
    Name      string    `gorm:"not null"`
    Value     string    `gorm:"not null"`
    Domain    string    `gorm:"not null"`
    Path      string    `gorm:"not null"`
    Expires   time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
} 
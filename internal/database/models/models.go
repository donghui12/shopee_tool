package models

import "time"

// Account 存储虾皮账号信息
type Account struct {
    ID        uint      `gorm:"primaryKey"`
    Username  string    `gorm:"uniqueIndex;not null"`
    Password  string    `gorm:"not null"`
    CreatedAt time.Time
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
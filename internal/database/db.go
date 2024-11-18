package database

import (
    "fmt"
    "gorm.io/gorm"
    "gorm.io/driver/mysql"
    "shopee_tool/internal/config"
    "shopee_tool/internal/database/models"
)

func InitDB(dbConfig *config.DatabaseConfig) (*gorm.DB, error) {
    db, err := gorm.Open(mysql.Open(dbConfig.GetDSN()), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("连接数据库失败: %v", err)
    }

    // 设置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
    sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)

    // 自动迁移数据库结构
    err = db.AutoMigrate(&models.Account{}, &models.Cookie{})
    if err != nil {
        return nil, fmt.Errorf("数据库迁移失败: %v", err)
    }

    return db, nil
} 
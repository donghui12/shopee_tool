package main

import (
    "log"
    "path/filepath"
    "github.com/gin-gonic/gin"
    "shopee_tool/internal/config"
    "shopee_tool/internal/database"
    "shopee_tool/internal/service"
    "shopee_tool/internal/api"
)

func main() {
    // 加载配置文件
    configPath := filepath.Join("configs", "dev.conf")
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("加载配置文件失败: %v", err)
    }

    // 初始化数据库连接
    db, err := database.InitDB(&cfg.Database)
    if err != nil {
        log.Fatalf("初始化数据库失败: %v", err)
    }

    // 创建登录服务
    loginService := service.NewLoginService(db)
    accountService := service.NewAccountService(db)
	activeCodeService := service.NewActiveCodeService(db)
	orderService := service.NewOrderService(db)

    // 创建 Gin 引擎
    engine := gin.Default()

    // 设置路由
    router := api.NewRouter(loginService, accountService, activeCodeService, orderService)
    router.SetupRoutes(engine)

    // 启动服务器
    if err := engine.Run(":8080"); err != nil {
        log.Fatalf("启动服务器失败: %v", err)
    }
}

package main

import (
    "log"
    "path/filepath"

    "github.com/gin-gonic/gin"
    "go.uber.org/zap"

    "shopee_tool/internal/config"
    "shopee_tool/internal/database"
    "shopee_tool/internal/service"
    "shopee_tool/internal/api"
    "shopee_tool/pkg/pool"
    "shopee_tool/pkg/constant"
    "shopee_tool/pkg/logger"
)

func main() {
    // 加载配置文件
    configPath := filepath.Join("configs", "dev.conf")
    cfg, err := config.LoadConfig(configPath)
    if err != nil {
        log.Fatalf("加载配置文件失败: %v", err)
    }
	// 初始化全局 WorkerPool
	pool.InitPool(constant.WorkerPoolSize)

    logger.Info("Worker pool initialized",
        zap.Int("worker_count", constant.WorkerPoolSize),
    )

    // 初始化数据库连接
    db, err := database.InitDB(&cfg.Database)
    if err != nil {
        log.Fatalf("初始化数据库失败: %v", err)
    }
    logger.Info("Database initialized",
        zap.String("database", cfg.Database.DBName),
    )

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

    logger.Info("Starting HTTP server",
        zap.String("address", ":8080"),
    )

    // 启动服务器
    if err := engine.Run(":8080"); err != nil {
        logger.Fatal("启动服务器失败", zap.Error(err))
	}
}

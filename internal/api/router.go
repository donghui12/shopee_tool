package api

import (
	"github.com/gin-gonic/gin"
	"shopee_tool/internal/service"
)

type Router struct {
	loginService *service.LoginService
	accountService *service.AccountService
	activeCodeService *service.ActiveCodeService
	orderService *service.OrderService
}

func NewRouter(loginService *service.LoginService, accountService *service.AccountService, 
	 activeCodeService *service.ActiveCodeService, orderService *service.OrderService) *Router {
	return &Router{
		loginService: loginService,
		accountService: accountService,
		activeCodeService: activeCodeService,
		orderService: orderService,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	{
		shopee := v1.Group("/shopee")
		{
			// 登录
			shopee.POST("/login", r.handleLogin)
			// 获取机器码
			shopee.GET("/mechine_code", r.handleGetMachineCode)
			// 更新账户机器码
			shopee.POST("/mechine_code", r.handleUpdateMachineCode)
			// 绑定账户和激活码
			shopee.POST("/bind_active_code", r.handleBindActiveCode)

			// 验证激活码
			shopee.GET("/verify_active_code", r.handleVerifyActiveCode)

			// 获取账户激活码
			shopee.GET("/active_code", r.handleGetActiveCode)
			// 创建激活码
			shopee.POST("/active_code", r.handleCreateActiveCode)
			// 更新账户下全部商品的库存
			shopee.POST("/update_order", r.handleUpdateOrder)
		}
	}
} 
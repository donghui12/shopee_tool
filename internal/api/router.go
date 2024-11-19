package api

import (
	"github.com/gin-gonic/gin"
	"shopee_tool/internal/service"
)

type Router struct {
	loginService *service.LoginService
	accountService *service.AccountService
	activeCodeService *service.ActiveCodeService
}

func NewRouter(loginService *service.LoginService, accountService *service.AccountService, 
	 activeCodeService *service.ActiveCodeService) *Router {
	return &Router{
		loginService: loginService,
		accountService: accountService,
		activeCodeService: activeCodeService,
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
			shopee.GET("/machine-code", r.handleGetMachineCode)
			// 更新账户机器码
			shopee.POST("/machine-code", r.handleUpdateMachineCode)
			// 获取账户激活码
			shopee.GET("/active-code", r.handleGetActiveCode)
			// 创建激活码
			shopee.POST("/active-code", r.handleCreateActiveCode)
		}
	}
} 
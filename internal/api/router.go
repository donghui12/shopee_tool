package api

import (
	"github.com/gin-gonic/gin"
	"shopee_tool/internal/service"
)

type Router struct {
	loginService *service.LoginService
	accountService *service.AccountService
}

func NewRouter(loginService *service.LoginService, accountService *service.AccountService) *Router {
	return &Router{
		loginService: loginService,
		accountService: accountService,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	{
		shopee := v1.Group("/shopee")
		{
			shopee.POST("/login", r.handleLogin)
			shopee.GET("/machine-code", r.handleGetMachineCode)
			shopee.POST("/machine-code", r.handleUpdateMachineCode)
		}
	}
} 
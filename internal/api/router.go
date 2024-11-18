package api

import (
	"github.com/gin-gonic/gin"
	"shopee_tool/internal/auth"
)

type Router struct {
	loginService *auth.LoginService
}

func NewRouter(loginService *auth.LoginService) *Router {
	return &Router{
		loginService: loginService,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	v1 := engine.Group("/v1")
	{
		shopee := v1.Group("/shopee")
		{
			shopee.POST("/login", r.handleLogin)
		}
	}
} 
package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Vcode    string `json:"vcode"`
}

type LoginResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Router) handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}

	// 调用登录服务
	cookies, err := r.loginService.Login(req.Username, req.Password, req.Vcode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, LoginResponse{
			Code:    500,
			Message: "登录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Code:    200,
		Message: "登录成功",
		Data: map[string]interface{}{
			"cookies": cookies,
		},
	})
}

type UpdateMachineCodeRequest struct {
	Username  string `json:"username"`
	MachineCode string `json:"machine_code"`
}

type UpdateMachineCodeResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (r *Router) handleUpdateMachineCode(c *gin.Context) {
	var req UpdateMachineCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateMachineCodeResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
	}
	if err := r.accountService.UpdateMachineCode(req.Username, req.MachineCode); err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "更新机器码失败: " + err.Error(),
		})
	}
}

func (r *Router) handleGetMachineCode(c *gin.Context) {
	username := c.Query("username")
	machineCode, err := r.accountService.GetMachineCode(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取机器码失败: " + err.Error(),
		})
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "获取机器码成功",
		Data:    machineCode,
	})
}

func (r *Router) handleGetActiveCode(c *gin.Context) {
	username := c.Query("username")
	activeCode, err := r.accountService.GetActiveCode(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取激活码失败: " + err.Error(),
		})
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "获取激活码成功",
		Data:    activeCode,
	})
}

func (r *Router) handleCreateActiveCode(c *gin.Context) {
	expiredAt := c.Query("expired_at")
	activeCode := uuid.New().String()
	activeCode, err := r.activeCodeService.CreateActiveCode(activeCode, expiredAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "创建激活码失败: " + err.Error(),
		})
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "创建激活码成功",
		Data:    activeCode,
	})
}

package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"github.com/google/uuid"
	"fmt"
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
	err := r.loginService.Login(req.Username, req.Password, req.Vcode)
	if err != nil {
		c.JSON(http.StatusOK, LoginResponse{
			Code:    410,
			Message: "登录失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Code:    200,
		Message: "登录成功",
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
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "更新机器码成功",
	})
}

func (r *Router) handleGetMachineCode(c *gin.Context) {
	username := c.Query("username")
	machineCode := c.Query("machine_code")
	fmt.Println("machineCode: ", machineCode)
	err := r.accountService.GetMachineCode(username, machineCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取机器码失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "获取机器码成功",
	})
}

func (r *Router) handleGetActiveCode(c *gin.Context) {
	username := c.Query("username")
	// 获取激活码
	activeCode, err := r.accountService.GetActiveCode(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取激活码失败: " + err.Error(),
		})
		return
	}
	// 验证激活码是否有效
	day, err := r.activeCodeService.GetActiveCode(activeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取激活码失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "获取激活码成功",
		Data:    day,
	})
}

type UpdateActiveCodeRequest struct {
	ExpiredAt string `json:"expired_at"`
}

func (r *Router) handleCreateActiveCode(c *gin.Context) {
	var req UpdateActiveCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateMachineCodeResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
	}

	currentActiveCode := uuid.New().String()
	activeCode, err := r.activeCodeService.CreateActiveCode(currentActiveCode, req.ExpiredAt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "创建激活码失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "创建激活码成功",
		Data:    activeCode,
	})
}

type BindActiveCodeRequest struct {
	Username  string `json:"username"`
	ActiveCode string `json:"active_code"`
}

// 绑定激活码
func (r *Router) handleBindActiveCode(c *gin.Context) {
	var req BindActiveCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateMachineCodeResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
	}
	// 绑定激活码
	err := r.accountService.UpdateActiveCode(req.Username, req.ActiveCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "绑定激活码失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "绑定激活码成功",
	})
}

type VerifyActiveCodeRequest struct {
	ActiveCode string `json:"active_code"`
}

func (r *Router) handleVerifyActiveCode(c *gin.Context) {
	activeCode := c.Query("active_code")
	// 验证激活码是否有效
	_, err := r.activeCodeService.GetActiveCode(activeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
				Code:    500,
				Message: "获取激活码失败: " + err.Error(),
			})
			return
	}
	err = r.accountService.GetActiveCodeByActiveCode(activeCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "激活码已被绑定: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "激活码有效",
	})
}

type UpdateOrderRequest struct {
	Username string `json:"username"`
	Days     int    `json:"days"`
}

// 更新账户下全部商品的库存
func (r *Router) handleUpdateOrder(c *gin.Context) {
	var req UpdateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, UpdateMachineCodeResponse{
			Code:    400,
			Message: "参数错误: " + err.Error(),
		})
		return
	}
	// 获取 cookies
	cookies, err := r.accountService.GetCookies(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "获取 cookies 失败: " + err.Error(),
		})
		return
	}

	// 更新账户下全部商品的库存
	err = r.orderService.UpdateOrder(cookies, req.Days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, UpdateMachineCodeResponse{
			Code:    500,
			Message: "更新库存失败: " + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, UpdateMachineCodeResponse{
		Code:    200,
		Message: "更新库存成功",
	})
}

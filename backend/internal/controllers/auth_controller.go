package controllers

import (
	"gobox/backend/internal/middleware"
	"gobox/backend/internal/services"
	"gobox/backend/pkg/response"
	validatorpkg "gobox/backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
	return &AuthController{service: service}
}

func (a *AuthController) SendRegisterCode(c *gin.Context) {
	var input services.SendRegisterCodePayload
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, 400, "请求参数错误", err.Error())
		return
	}
	if err := validatorpkg.Validate(input); err != nil {
		response.Fail(c, 422, "参数校验失败", err.Error())
		return
	}

	result, err := a.service.SendRegisterCode(input)
	if err != nil {
		response.Fail(c, 400, "发送验证码失败", err.Error())
		return
	}

	message := "验证码已发送，请查收邮箱"
	if result.PreviewCode != "" {
		message = "开发模式未配置邮件服务，已返回预览验证码"
	}
	response.OK(c, message, result)
}

func (a *AuthController) Register(c *gin.Context) {
	var input services.RegisterPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, 400, "请求参数错误", err.Error())
		return
	}
	if err := validatorpkg.Validate(input); err != nil {
		response.Fail(c, 422, "参数校验失败", err.Error())
		return
	}

	user, token, err := a.service.Register(input)
	if err != nil {
		response.Fail(c, 400, "注册失败", err.Error())
		return
	}
	response.Created(c, "注册成功", gin.H{"user": user, "token": token})
}

func (a *AuthController) Login(c *gin.Context) {
	var input services.LoginPayload
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, 400, "请求参数错误", err.Error())
		return
	}
	if err := validatorpkg.Validate(input); err != nil {
		response.Fail(c, 422, "参数校验失败", err.Error())
		return
	}

	user, token, err := a.service.Login(input)
	if err != nil {
		response.Fail(c, 401, "登录失败", err.Error())
		return
	}
	response.OK(c, "登录成功", gin.H{"user": user, "token": token})
}

func (a *AuthController) Profile(c *gin.Context) {
	current, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未授权", "missing user context")
		return
	}
	user, err := a.service.Profile(current.ID)
	if err != nil {
		response.Fail(c, 404, "用户不存在", err.Error())
		return
	}
	response.OK(c, "获取成功", user)
}

func (a *AuthController) SavePreferences(c *gin.Context) {
	current, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未授权", "missing user context")
		return
	}
	var input struct {
		Preferences map[string]string `json:"preferences"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, 400, "请求参数错误", err.Error())
		return
	}
	if err := a.service.SavePreferences(current.ID, input.Preferences); err != nil {
		response.Fail(c, 500, "保存失败", err.Error())
		return
	}
	response.OK(c, "保存成功", input.Preferences)
}

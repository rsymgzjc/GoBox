package controllers

import (
	"gobox/backend/internal/middleware"
	"gobox/backend/internal/services"
	"gobox/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ToolController struct {
	service *services.ToolService
}

func NewToolController(service *services.ToolService) *ToolController {
	return &ToolController{service: service}
}

func (t *ToolController) List(c *gin.Context) {
	tools, err := t.service.List(c.Query("category"))
	if err != nil {
		response.Fail(c, 500, "获取工具失败", err.Error())
		return
	}
	response.OK(c, "获取成功", tools)
}

func (t *ToolController) Run(c *gin.Context) {
	var input services.RunToolInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, 400, "请求参数错误", err.Error())
		return
	}

	var userID *uint
	if current, ok := middleware.CurrentUser(c); ok {
		userID = &current.ID
	}

	result, err := t.service.Run(c.Param("slug"), input, userID)
	if err != nil {
		response.Fail(c, 400, "工具执行失败", err.Error())
		return
	}
	response.OK(c, "执行成功", result)
}

func (t *ToolController) Summary(c *gin.Context) {
	summary, err := t.service.Summary()
	if err != nil {
		response.Fail(c, 500, "获取统计失败", err.Error())
		return
	}
	response.OK(c, "获取成功", summary)
}

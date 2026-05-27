package controllers

import (
	"time"

	"gobox/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func Health(c *gin.Context) {
	response.OK(c, "ok", gin.H{
		"status":    "up",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}

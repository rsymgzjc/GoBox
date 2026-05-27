package response

import "github.com/gin-gonic/gin"

type Envelope struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(200, Envelope{Success: true, Message: message, Data: data})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(201, Envelope{Success: true, Message: message, Data: data})
}

func Fail(c *gin.Context, status int, message, err string) {
	c.JSON(status, Envelope{Success: false, Message: message, Error: err})
}

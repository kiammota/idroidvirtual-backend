package gateways

import "github.com/gin-gonic/gin"

type Response struct {
	Data    any    `json:"data"`
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func JSON(c *gin.Context, status int, data any, success bool, message string) {
	c.JSON(status, Response{
		Data:    data,
		Success: success,
		Message: message,
	})
}

func Success(c *gin.Context, status int, data any, message string) {
	JSON(c, status, data, true, message)
}

func Error(c *gin.Context, status int, message string) {
	JSON(c, status, nil, false, message)
}

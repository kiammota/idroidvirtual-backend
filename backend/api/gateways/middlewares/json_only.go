package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func JsonOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Métodos que normalmente possuem payload
		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {

			contentType := c.GetHeader("Content-Type")

			if !strings.HasPrefix(contentType, "application/json") {
				c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{
					"data":    nil,
					"success": false,
					"message": "Content-Type deve ser application/json",
				})
				return
			}
		}

		c.Next()
	}
}

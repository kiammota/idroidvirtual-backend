package middlewares

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

func JsonNotEmptyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		if c.Request.Method == http.MethodPost ||
			c.Request.Method == http.MethodPut ||
			c.Request.Method == http.MethodPatch {

			body, err := io.ReadAll(c.Request.Body)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Não foi possível ler o corpo da requisição",
				})
				return
			}

			// Recoloca o body para o Gin usar depois
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

			// Body vazio
			if len(bytes.TrimSpace(body)) == 0 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "O corpo da requisição não pode ser vazio",
				})
				return
			}

			var data interface{}

			err = json.Unmarshal(body, &data)

			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "JSON inválido",
				})
				return
			}

			// JSON {} ou todos campos nulos
			if isEmptyJSON(data) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Os dados enviados não podem ser vazios",
				})
				return
			}
		}

		c.Next()
	}
}

func isEmptyJSON(data interface{}) bool {

	if data == nil {
		return true
	}

	value := reflect.ValueOf(data)

	switch value.Kind() {

	case reflect.Map:
		if value.Len() == 0 {
			return true
		}

		for _, key := range value.MapKeys() {
			item := value.MapIndex(key)

			if item.Interface() != nil {
				return false
			}
		}

		return true

	case reflect.Slice:
		return value.Len() == 0

	default:
		return false
	}
}

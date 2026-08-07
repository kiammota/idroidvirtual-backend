package middlewares

import (
	"bytes"
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

func LoggerBodyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Se houver corpo na requisição (POST, PUT, PATCH, etc.)
		if c.Request.Body != nil {
			// 1. Lê todo o conteúdo do body
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				log.Printf("[LOG] Erro ao ler body: %v\n", err)
			} else {
				// 2. Imprime no console a Rota, Método e o JSON bruto recebido
				log.Printf("\n--- [REQUISIÇÃO %s %s] ---", c.Request.Method, c.Request.URL.Path)
				log.Printf("[LOG BODY]: %s\n---------------------------------------\n", string(bodyBytes))

				// 3. RESTAURA o body para que os Handlers (ShouldBindJSON) consigam ler novamente
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			}
		}

		c.Next()
	}
}

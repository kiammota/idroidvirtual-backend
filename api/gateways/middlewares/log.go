package middlewares

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var errorLogger = log.New(
	getLogFile(),
	"",
	log.Ldate|log.Ltime,
)

func getLogFile() *os.File {
	file, err := os.OpenFile(
		"logs/error.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)

	if err != nil {
		panic(err)
	}

	return file
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		inicio := time.Now()

		c.Next()

		status := c.Writer.Status()

		log.Printf(
			"%s %s | %d | %s",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			time.Since(inicio),
		)

		if status >= 500 {
			errorLogger.Printf(
				`
Erro 500 detectado
Data: %s
Método: %s
Rota: %s
IP: %s
Status: %d
`,
				time.Now().Format("2006-01-02 15:04:05"),
				c.Request.Method,
				c.Request.URL.Path,
				c.ClientIP(),
				status,
			)
		}
	}
}

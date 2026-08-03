package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var JWTSecret = []byte("sua-chave-secreta")

func AuthJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token não informado",
			})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")

		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token inválido",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, jwt.ErrSignatureInvalid
			}

			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token expirado ou inválido",
			})
			ctx.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)

		if ok {

			if usuarioID, exists := claims["usuario_id"]; exists {
				ctx.Set("usuario_id", usuarioID)
			}

			if email, exists := claims["email"]; exists {
				ctx.Set("email", email)
			}
		}

		ctx.Next()
	}
}

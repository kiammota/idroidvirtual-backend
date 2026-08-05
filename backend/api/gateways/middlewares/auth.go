package middlewares

import (
	"net/http"
	"strings"

	"backend/services"

	"github.com/gin-gonic/gin"
)

// AuthJWT recebe a instância do TokenService criada no bootstrap do servidor.
func AuthJWT(tokenService *services.TokenService) gin.HandlerFunc {
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
				"message": "Formato de token inválido",
			})
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		// Usa o TokenService para validar centralizadamente
		claims, err := tokenService.ValidarToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Token expirado ou inválido",
			})
			ctx.Abort()
			return
		}

		// Grava os dados da struct no contexto do Gin
		ctx.Set("usuario_id", claims.UsuarioID)
		ctx.Set("email", claims.Email)

		ctx.Next()
	}
}

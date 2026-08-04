package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(
	router *gin.RouterGroup,
	authController *controllers.AuthController,
) {
	auth := router.Group("/auth")
	{
		auth.POST("/login", authController.UsuarioLogar)
		auth.POST("/elevate", authController.UsuarioSenha)
	}
}

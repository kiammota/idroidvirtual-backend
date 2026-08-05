package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func UsuarioRoutes(
	router *gin.RouterGroup,
	usuarioController *controllers.UsuarioController,
) {
	basePath := "/users"

	router.POST(
		basePath,
		usuarioController.CriarUsuario,
	)

	router.GET(
		basePath+"/:id",
		usuarioController.BuscarUsuario,
	)

	router.POST(
		basePath+"/elevate",
		usuarioController.UsuarioSenha,
	)
}

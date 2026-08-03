package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func UsuarioRoutes(
	router *gin.RouterGroup,
	usuarioController *controllers.UsuarioController,
) {

	router.POST(
		"/usuarios",
		usuarioController.CriarUsuario,
	)

	router.GET(
		"/usuarios/:id",
		usuarioController.BuscarUsuario,
	)

}

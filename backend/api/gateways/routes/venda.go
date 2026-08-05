package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func VendaRoutes(
	router *gin.RouterGroup,
	vendaController *controllers.VendaController,
) {
	vendas := router.Group("/vendas")
	{
		vendas.POST("", vendaController.Vender)
		vendas.POST("/:id/cancelar", vendaController.Cancelar)
	}
}

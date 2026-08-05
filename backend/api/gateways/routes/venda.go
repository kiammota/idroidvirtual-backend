package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func VendaRoutes(router *gin.RouterGroup, salesController *controllers.VendasController) {
	salesGroup := router.Group("/sales")
	{
		// Rotas fixas/específicas vem PRIMEIRO para não conflitar com rotas dinâmicas como /:id
		salesGroup.GET("/dia", salesController.GetSalesOfDay)
		salesGroup.GET("/mes", salesController.GetSalesOfMonth)

		// Rota para criação de venda
		salesGroup.POST("", salesController.Create)

		// Rota para cancelamento de venda por ID
		salesGroup.DELETE("/:id", salesController.Cancel)
	}
}

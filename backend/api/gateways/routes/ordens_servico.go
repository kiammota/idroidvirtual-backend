package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func OrdensServicoRoutes(r *gin.RouterGroup, ctrl *controllers.OrdemServicoController) {

	serviceOrders := r.Group("/service-orders")

	serviceOrders.POST("", ctrl.Criar)
	serviceOrders.GET("", ctrl.Listar)

	serviceOrders.GET("/three-months", ctrl.UltimosTresMeses)
	serviceOrders.GET("/six-months", ctrl.UltimosSeisMeses)

	// serviceOrders.GET("/client-id/:cliente_id", ctrl.BuscarPorClienteID)
	// serviceOrders.GET("/client-name/:client_name", ctrl.BuscarPorNomeCliente)
	// serviceOrders.GET("/client-email/:client_email", ctrl.BuscarPorEmailCliente)
	// serviceOrders.GET("/client-number/:client_number", ctrl.BuscarPorNumeroCliente)

	serviceOrders.GET("/:id", ctrl.BuscarPorID)
	serviceOrders.PATCH("/:id", ctrl.Atualizar)
	serviceOrders.DELETE("/:id", ctrl.Deletar)
}

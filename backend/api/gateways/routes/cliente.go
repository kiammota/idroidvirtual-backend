package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

// ClienteRoutes registra o RouterGroup de clientes na aplicação Gin
func ClienteRoutes(router *gin.RouterGroup, clientController *controllers.ClientController) {
	clientsGroup := router.Group("/clients")
	{
		// Criar cliente
		clientsGroup.POST("/", clientController.Create)

		// Listar todos os clientes (aceita ?limit=10&offset=0)
		clientsGroup.GET("/", clientController.ListAll)

		// Buscas específicas por query params (devem vir antes de /:id para evitar conflitos de rota)
		clientsGroup.GET("/number", clientController.SearchByPhone)
		clientsGroup.GET("/address", clientController.SearchByAddress)
		clientsGroup.GET("/name", clientController.SearchByName)
		clientsGroup.GET("/email", clientController.SearchByEmail)

		// Buscar cliente por ID
		clientsGroup.GET("/:id", clientController.GetByID)
	}
}

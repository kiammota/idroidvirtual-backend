package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func RegistrarRotasTroca(
	router *gin.RouterGroup,
	trocaHandler *controllers.TrocaController,
) {

	trocas := router.Group("/trocas")
	{
		// Criar troca vinculada a uma venda
		trocas.POST("", trocaHandler.Criar)

		// Listar todas as trocas (paginado)
		trocas.GET("", trocaHandler.Listar)

		// Listar trocas por períodos (paginados)
		trocas.GET("/ultimo-mes", trocaHandler.ListarUltimoMes)
		trocas.GET("/tres-meses", trocaHandler.ListarUltimosTresMeses)
		trocas.GET("/seis-meses", trocaHandler.ListarUltimosSeisMeses)

		// Buscar troca de uma venda específica
		trocas.GET("/venda/:vendaId", trocaHandler.ObterPorVendaID)
	}
}

package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func RegistrarRotasTroca(
	router *gin.RouterGroup,
	trocaHandler *controllers.TrocaController,
) {

	v1 := router.Group("/api/v1")
	{
		trocas := v1.Group("/trocas")
		{
			// Criar troca vinculada a uma venda
			trocas.POST(
				"",
				trocaHandler.Criar,
			)

			// Buscar troca de uma venda específica
			trocas.GET(
				"/venda/:vendaId",
				trocaHandler.ObterPorVendaID,
			)
		}
	}
}

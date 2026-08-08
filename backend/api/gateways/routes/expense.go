package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func RegistrarDespesas(router *gin.RouterGroup, c controllers.DespesaController) {
	despesas := router.Group("/expense")
	{
		despesas.POST("", c.Criar)
		despesas.GET("", c.Listar)
		despesas.GET("/:id", c.BuscarPorID)
		despesas.PUT("/:id", c.Atualizar)
		despesas.DELETE("/:id", c.Deletar)
	}
}

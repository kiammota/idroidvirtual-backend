package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func FinanceiroRoutes(router *gin.RouterGroup, c controllers.FinanceiroController) {
	financeiro := router.Group("/financial")
	{
		financeiro.GET("/today", c.ObterHoje)
		financeiro.GET("/month", c.ObterMesAtual)
		financeiro.GET("/period", c.ObterPorIntervalo)
	}
}

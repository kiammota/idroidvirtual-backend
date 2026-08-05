package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

// ClosingRoutes registra o RouterGroup de fechamento na aplicação Gin
func ClosingRoutes(router *gin.RouterGroup, closingController *controllers.ClosingController) {
	closingGroup := router.Group("/closing")
	{
		// Fechamento do dia
		closingGroup.GET("/dia", closingController.GetDayClosing)

		// Fechamento do mês
		closingGroup.GET("/mes", closingController.GetMonthClosing)
	}
}

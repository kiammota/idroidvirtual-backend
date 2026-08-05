package routes

import (
	"backend/gateways/controllers"

	"github.com/gin-gonic/gin"
)

func AuditoriaRoute(h *gin.RouterGroup, c *controllers.AuditoriaController) {
	h.GET("/auditoria", c.ListarAuditorias)
}

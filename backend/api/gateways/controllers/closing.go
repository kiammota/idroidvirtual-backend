package controllers

import (
	"net/http"

	"backend/gateways"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type ClosingController struct {
	vendaService *services.VendaService
}

func NewClosingController(vendaService *services.VendaService) *ClosingController {
	return &ClosingController{
		vendaService: vendaService,
	}
}

// GET /closing/dia
func (c *ClosingController) GetDayClosing(ctx *gin.Context) {
	fechamento, err := c.vendaService.FechamentoDoDia(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, "Erro ao gerar fechamento do dia: "+err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, fechamento, "Fechamento do dia gerado com sucesso")
}

// GET /closing/mes
func (c *ClosingController) GetMonthClosing(ctx *gin.Context) {
	fechamento, err := c.vendaService.FechamentoDoMes(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, "Erro ao gerar fechamento do mês: "+err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, fechamento, "Fechamento do mês gerado com sucesso")
}

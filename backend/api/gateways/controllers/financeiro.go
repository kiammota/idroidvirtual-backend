package controllers

import (
	"net/http"
	"time"

	"backend/gateways"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type FinanceiroController struct {
	financeiroService *services.FinanceiroService
}

func NewFinanceiroController(financeiroService *services.FinanceiroService) *FinanceiroController {
	return &FinanceiroController{
		financeiroService: financeiroService,
	}
}

// ObterHoje retorna o resumo financeiro do dia atual
func (c *FinanceiroController) ObterHoje(ctx *gin.Context) {
	resumo, err := c.financeiroService.FechamentoPorPeriodo(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, "Erro ao buscar resumo financeiro do dia: "+err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, resumo, "Resumo financeiro do dia obtido com sucesso")
}

// ObterMesAtual retorna o resumo financeiro do mês corrente
func (c *FinanceiroController) ObterMesAtual(ctx *gin.Context) {
	resumo, err := c.financeiroService.ObterFinanceiroDoMes(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, "Erro ao buscar resumo financeiro do mês: "+err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, resumo, "Resumo financeiro do mês obtido com sucesso")
}

// ObterPorIntervalo busca o resumo financeiro entre duas datas (Formato: YYYY-MM-DD)
func (c *FinanceiroController) ObterPorIntervalo(ctx *gin.Context) {
	dataInicioStr := ctx.Query("inicio")
	dataFimStr := ctx.Query("fim")

	if dataInicioStr == "" || dataFimStr == "" {
		gateways.Error(ctx, http.StatusBadRequest, "Os parâmetros 'inicio' e 'fim' são obrigatórios. Ex: ?inicio=2026-08-01&fim=2026-08-31")
		return
	}

	inicio, err := time.Parse("2006-01-02", dataInicioStr)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Formato de data inicial inválido. Use AAAA-MM-DD")
		return
	}

	fimRaw, err := time.Parse("2006-01-02", dataFimStr)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Formato de data final inválido. Use AAAA-MM-DD")
		return
	}

	// Ajusta o fim do dia para incluir todas as horas até 23:59:59
	fim := time.Date(fimRaw.Year(), fimRaw.Month(), fimRaw.Day(), 23, 59, 59, 999999999, fimRaw.Location())

	resumo, err := c.financeiroService.ObterFinanceiroPorIntervalo(ctx.Request.Context(), inicio, fim)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, resumo, "Resumo financeiro do período obtido com sucesso")
}

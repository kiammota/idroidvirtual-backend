package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type VendasController struct {
	vendaService *services.VendaService
}

// Construtor limpo: recebe APENAS a service
func NewVendasController(vendaService *services.VendaService) *VendasController {
	return &VendasController{
		vendaService: vendaService,
	}
}

type CancelSaleInput struct {
	UsuarioID int64  `json:"usuario_id" binding:"required"`
	Motivo    string `json:"motivo"`
}

// POST /sales
func (c *VendasController) Create(ctx *gin.Context) {
	var input models.Venda

	if err := ctx.ShouldBindJSON(&input); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos para registro da venda")
		return
	}

	venda, err := c.vendaService.ExecutarVenda(
		ctx.Request.Context(),
		input.UsuarioID,
		input.ProdutoID,
		input.Quantidade,
		input.TipoPagamento,
		input.ClienteID,
	)

	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, venda, "Venda realizada com sucesso")
}

// DELETE /sales/:id
func (c *VendasController) Cancel(ctx *gin.Context) {
	idParam := ctx.Param("id")
	vendaID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || vendaID <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID da venda é inválido")
		return
	}

	var input CancelSaleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "O campo usuario_id é obrigatório para cancelar")
		return
	}

	err = c.vendaService.CancelarVenda(
		ctx.Request.Context(),
		input.UsuarioID,
		vendaID,
		input.Motivo,
	)

	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Venda cancelada com sucesso")
}

// GET /sales/dia
func (c *VendasController) GetSalesOfDay(ctx *gin.Context) {
	vendas, err := c.vendaService.ObterVendasDoDia(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, vendas, "Vendas do dia recuperadas com sucesso")
}

// GET /sales/mes
func (c *VendasController) GetSalesOfMonth(ctx *gin.Context) {
	vendas, err := c.vendaService.ObterVendasDoMes(ctx.Request.Context())
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, vendas, "Vendas do mês recuperadas com sucesso")
}

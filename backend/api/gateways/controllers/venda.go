package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type VendaController struct {
	VendaService *services.VendaService
}

func NewVendaController(vendaService *services.VendaService) *VendaController {
	return &VendaController{
		VendaService: vendaService,
	}
}

type ExecutarVendaInput struct {
	ProdutoID     int64  `json:"produto_id" binding:"required"`
	Quantidade    int    `json:"quantidade" binding:"required"`
	TipoPagamento string `json:"tipo_pagamento" binding:"required"`
}

type CancelarVendaInput struct {
	Motivo string `json:"motivo" binding:"required"`
}

// Vender processa a realização de uma nova venda
// POST /vendas
func (vc *VendaController) Vender(ctx *gin.Context) {
	var input ExecutarVendaInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos: "+err.Error())
		return
	}

	usuarioID := ctx.GetInt64("usuario_id")

	venda, err := vc.VendaService.ExecutarVenda(
		ctx.Request.Context(),
		usuarioID,
		input.ProdutoID,
		input.Quantidade,
		input.TipoPagamento,
	)

	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, venda, "Venda realizada com sucesso")
}

// Cancelar processa o estorno de uma venda existente
// POST /vendas/:id/cancelar
func (vc *VendaController) Cancelar(ctx *gin.Context) {
	vendaID, _ := strconv.ParseInt(ctx.Param("id"), 10, 64)

	var input CancelarVendaInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "O motivo do cancelamento é obrigatório")
		return
	}

	usuarioID := ctx.GetInt64("usuario_id")

	err := vc.VendaService.CancelarVenda(
		ctx.Request.Context(),
		usuarioID,
		vendaID,
		input.Motivo,
	)

	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Venda cancelada e estoque estornado com sucesso")
}

package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/gateways/helpers"
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

// POST /sales
func (c *VendasController) Create(ctx *gin.Context) {
	var input struct {
		ProdutoID     int64  `json:"produto_id"`
		Quantidade    int    `json:"quantidade"`
		TipoPagamento string `json:"tipo_pagamento"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos para registro da venda")
		return
	}

	// 1. Pega o empresa_id do token JWT usando o helper
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "Não foi possível identificar a empresa do usuário")
		return
	}

	// 2. Pega o usuario_id do token JWT usando o helper
	usuarioID, err := helpers.GetUserID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "Não foi possível identificar o usuário autenticado")
		return
	}

	// 3. Executa a venda passando o empresaID obtido com segurança do token
	venda, err := c.vendaService.ExecutarVenda(
		ctx.Request.Context(),
		empresaID, // <-- Aqui está o empresa_id correto vindo do token
		usuarioID,
		input.ProdutoID,
		input.Quantidade,
		input.TipoPagamento,
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

	// Recupera o usuario_id usando o helper padronizado
	usuarioID, err := helpers.GetUserID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "Usuário não autenticado ou ID inválido")
		return
	}

	// Executa a regra de negócio do serviço com os IDs validados
	err = c.vendaService.CancelarVenda(
		ctx.Request.Context(),
		usuarioID,
		vendaID,
	)

	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Venda cancelada com sucesso")
}

// GET /sales/dia
func (c *VendasController) GetSalesOfDay(ctx *gin.Context) {
	empresaId, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "Não foi possível pegar o ID da empresa!")
		return
	}

	vendas, err := c.vendaService.ObterVendasDoDia(ctx.Request.Context(), empresaId)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, vendas, "Vendas do dia recuperadas com sucesso")
}

// GET /sales/mes
func (c *VendasController) GetSalesOfMonth(ctx *gin.Context) {
	empresaId, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "Não foi possível pegar o ID da empresa!")
		return
	}

	vendas, err := c.vendaService.ObterVendasDoMes(ctx.Request.Context(), empresaId)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, vendas, "Vendas do mês recuperadas com sucesso")
}

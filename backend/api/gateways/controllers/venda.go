package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/gateways/helpers"
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

	// 1. Recupera o usuario_id injetado no contexto pelo Middleware de Auth
	usuarioIDValue, exists := ctx.Get("usuario_id")
	if !exists {
		gateways.Error(ctx, http.StatusUnauthorized, "Usuário não autenticado")
		return
	}

	// 2. Garante a conversão correta do tipo (int64)
	usuarioID, ok := usuarioIDValue.(int64)
	if !ok {
		// Trata caso o middleware tenha salvo como float64 (comum no unmarshal de JSON do JWT)
		if floatVal, okFloat := usuarioIDValue.(float64); okFloat {
			usuarioID = int64(floatVal)
		} else {
			gateways.Error(ctx, http.StatusInternalServerError, "Erro ao processar identificador do usuário")
			return
		}
	}

	// 3. Executa a regra de negócio do serviço com os IDs validados
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
		gateways.Error(ctx, 500, "não foi possível pegar o id da empreas!")
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
		gateways.Error(ctx, 500, "não foi possível pegar o id da empreas!")
	}
	vendas, err := c.vendaService.ObterVendasDoMes(ctx.Request.Context(), empresaId)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, vendas, "Vendas do mês recuperadas com sucesso")
}

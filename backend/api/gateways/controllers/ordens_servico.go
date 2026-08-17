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

type OrdemServicoController struct {
	Service *services.OrdensServicoService
}

func NewOrdemServicoController(service *services.OrdensServicoService) *OrdemServicoController {
	return &OrdemServicoController{Service: service}
}

// Helper interno para extrair paginação da URL (?limit=10&offset=0)
func parsePaginationParams(ctx *gin.Context) (int, int) {
	limit, err := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		offset = 0
	}

	return limit, offset
}

// POST /service-orders
func (c *OrdemServicoController) Criar(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado: "+err.Error())
		return
	}

	usuarioID, err := helpers.GetUserID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado: "+err.Error())
		return
	}

	var ordem models.OrdemServico
	if !gateways.BindJSON(ctx, &ordem) {
		return
	}

	// Injeta os IDs extraídos do token na ordem de serviço
	ordem.EmpresaID = empresaID
	ordem.UsuarioID = usuarioID

	// Executa a criação
	err = c.Service.Criar(ctx.Request.Context(), &ordem)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, ordem, "ordem de serviço criada com sucesso")
}

// GET /service-orders
func (c *OrdemServicoController) Listar(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado")
		return
	}

	limit, offset := parsePaginationParams(ctx)

	// Passa o empresaID para o service
	ordens, err := c.Service.Listar(ctx.Request.Context(), empresaID, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens listadas com sucesso")
}

// GET /service-orders/:id
func (c *OrdemServicoController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "id inválido")
		return
	}

	// Como definido no repository, BuscarPorID exige apenas o ID da OS
	ordem, err := c.Service.BuscarPorID(ctx.Request.Context(), id)
	if err != nil {
		gateways.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordem, "ordem de serviço encontrada")
}

// PATCH /service-orders/:id
func (c *OrdemServicoController) Atualizar(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado")
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "id inválido")
		return
	}

	var patch models.OrdemServico
	if !gateways.BindJSON(ctx, &patch) {
		return
	}

	// Garante que o usuário não mude a OS de empresa
	patch.EmpresaID = empresaID

	err = c.Service.Atualizar(ctx.Request.Context(), id, &patch)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "ordem de serviço atualizada com sucesso")
}

// DELETE /service-orders/:id
func (c *OrdemServicoController) Deletar(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "id inválido")
		return
	}

	// O Service deve receber o contexto e o ID. (O repository só precisa do ID).
	err = c.Service.Deletar(ctx.Request.Context(), id)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "ordem de serviço removida com sucesso")
}

// GET /service-orders/three-months
func (c *OrdemServicoController) UltimosTresMeses(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado")
		return
	}

	// Adiciona empresaID
	ordens, err := c.Service.UltimosTresMeses(ctx.Request.Context(), empresaID)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens dos últimos 3 meses obtidas com sucesso")
}

// GET /service-orders/six-months
func (c *OrdemServicoController) UltimosSeisMeses(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusUnauthorized, "não autorizado")
		return
	}

	// Adiciona empresaID
	ordens, err := c.Service.UltimosSeisMeses(ctx.Request.Context(), empresaID)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens dos últimos 6 meses obtidas com sucesso")
}

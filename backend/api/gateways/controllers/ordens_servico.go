package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
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
	var ordem models.OrdemServico

	if !gateways.BindJSON(ctx, &ordem) {
		return
	}

	err := c.Service.Criar(ctx, &ordem)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, ordem, "ordem de serviço criada com sucesso")
}

// GET /service-orders
func (c *OrdemServicoController) Listar(ctx *gin.Context) {
	limit, offset := parsePaginationParams(ctx)

	ordens, err := c.Service.Listar(ctx, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens listadas com sucesso")
}

func (c *OrdemServicoController) ListarPorStatus(ctx *gin.Context) {
	limit, offset := parsePaginationParams(ctx)
	status := ctx.Query("status")

	if status == "" {
		gateways.Error(ctx, http.StatusBadRequest, "o parâmetro status é obrigatório")
		return
	}

	ordens, err := c.Service.ListarPorStatus(ctx, status, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens listadas por status com sucesso")
}

// GET /service-orders/:id
func (c *OrdemServicoController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "id inválido")
		return
	}

	ordem, err := c.Service.BuscarPorID(ctx, id)
	if err != nil {
		gateways.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordem, "ordem de serviço encontrada")
}

// PATCH /service-orders/:id
func (c *OrdemServicoController) Atualizar(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "id inválido")
		return
	}
	var patch models.OrdemServico

	if !gateways.BindJSON(ctx, &patch) {
		return
	}

	err = c.Service.Atualizar(ctx, id, &patch)
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

	err = c.Service.Deletar(ctx, id)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "ordem de serviço removida com sucesso")
}

// GET /service-orders/client-id/:cliente_id
func (c *OrdemServicoController) BuscarPorClienteID(ctx *gin.Context) {
	clienteID, err := strconv.ParseInt(ctx.Param("cliente_id"), 10, 64)
	if err != nil || clienteID <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "cliente_id inválido")
		return
	}

	limit, offset := parsePaginationParams(ctx)

	ordens, err := c.Service.BuscarPorClienteID(ctx, clienteID, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens do cliente encontradas")
}

// GET /service-orders/client-name/:client_name
func (c *OrdemServicoController) BuscarPorNomeCliente(ctx *gin.Context) {
	nome := ctx.Param("client_name")
	limit, offset := parsePaginationParams(ctx)

	ordens, err := c.Service.BuscarPorNomeCliente(ctx, nome, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens buscadas por nome com sucesso")
}

// GET /service-orders/client-email/:client_email
func (c *OrdemServicoController) BuscarPorEmailCliente(ctx *gin.Context) {
	email := ctx.Param("client_email")
	limit, offset := parsePaginationParams(ctx)

	ordens, err := c.Service.BuscarPorEmailCliente(ctx, email, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens buscadas por email com sucesso")
}

// GET /service-orders/client-number/:client_number
func (c *OrdemServicoController) BuscarPorNumeroCliente(ctx *gin.Context) {
	numero := ctx.Param("client_number")
	limit, offset := parsePaginationParams(ctx)

	ordens, err := c.Service.BuscarPorNumeroCliente(ctx, numero, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens buscadas por número de telefone com sucesso")
}

// GET /service-orders/three-months
func (c *OrdemServicoController) UltimosTresMeses(ctx *gin.Context) {
	ordens, err := c.Service.UltimosTresMeses(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens dos últimos 3 meses obtidas com sucesso")
}

// GET /service-orders/six-months
func (c *OrdemServicoController) UltimosSeisMeses(ctx *gin.Context) {
	ordens, err := c.Service.UltimosSeisMeses(ctx)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, ordens, "ordens dos últimos 6 meses obtidas com sucesso")
}

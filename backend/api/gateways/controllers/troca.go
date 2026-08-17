package controllers

import (
	"context"
	"net/http"
	"strconv"

	"backend/gateways/helpers"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type TrocaController struct {
	trocaService services.TrocaService
}

func NewTrocaController(service services.TrocaService) *TrocaController {
	return &TrocaController{
		trocaService: service,
	}
}

// POST /trocas
func (h *TrocaController) Criar(c *gin.Context) {
	// Utilizando o helper para pegar o ID da empresa do token
	empresaID, err := helpers.GetEmpresaID(c)
	if err != nil || empresaID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Não foi possível identificar a empresa na requisição",
		})
		return
	}

	var input models.Troca
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	// Injeta o ID da empresa recuperado do token
	input.EmpresaID = empresaID

	if err := h.trocaService.ProcessarTroca(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /trocas/:id
func (h *TrocaController) ObterPorID(c *gin.Context) {
	// Utilizando o helper
	empresaID, err := helpers.GetEmpresaID(c)
	if err != nil || empresaID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Não foi possível identificar a empresa na requisição",
		})
		return
	}

	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID inválido",
		})
		return
	}

	troca, err := h.trocaService.ObterPorID(
		c.Request.Context(),
		id,
		empresaID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if troca == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Troca não encontrada",
		})
		return
	}

	c.JSON(http.StatusOK, troca)
}

// GET /trocas?limit=10&offset=0
func (h *TrocaController) Listar(c *gin.Context) {
	h.executarListagem(c, h.trocaService.Listar)
}

// GET /trocas/ultimo-mes?limit=10&offset=0
func (h *TrocaController) ListarUltimoMes(c *gin.Context) {
	h.executarListagem(c, h.trocaService.ListarUltimoMes)
}

// GET /trocas/tres-meses?limit=10&offset=0
func (h *TrocaController) ListarUltimosTresMeses(c *gin.Context) {
	h.executarListagem(c, h.trocaService.ListarUltimosTresMeses)
}

// GET /trocas/seis-meses?limit=10&offset=0
func (h *TrocaController) ListarUltimosSeisMeses(c *gin.Context) {
	h.executarListagem(c, h.trocaService.ListarUltimosSeisMeses)
}

// Função auxiliar para evitar repetição de código nas chamadas de listagem
func (h *TrocaController) executarListagem(
	c *gin.Context,
	fn func(ctx context.Context, empresaID int64, limit, offset int) ([]*models.Troca, error),
) {
	// Utilizando o helper aqui também
	empresaID, err := helpers.GetEmpresaID(c)
	if err != nil || empresaID <= 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Não foi possível identificar a empresa na requisição",
		})
		return
	}

	limit, offset := parsePaginacao(c)

	trocas, err := fn(c.Request.Context(), empresaID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trocas)
}

// Função auxiliar para extrair e validar limit e offset dos query params
func parsePaginacao(c *gin.Context) (limit, offset int) {
	limit, _ = strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ = strconv.Atoi(c.DefaultQuery("offset", "0"))

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

package controllers

import (
	"net/http"
	"strconv"

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

// POST /api/v1/trocas
func (h *TrocaController) Criar(c *gin.Context) {
	var input models.Troca

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dados inválidos: " + err.Error(),
		})
		return
	}

	if err := h.trocaService.ProcessarTroca(c.Request.Context(), &input); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, input)
}

// GET /api/v1/trocas/venda/:vendaId
func (h *TrocaController) ObterPorVendaID(c *gin.Context) {
	vendaID, err := strconv.ParseInt(c.Param("vendaId"), 10, 64)

	if err != nil || vendaID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "ID da venda inválido",
		})
		return
	}

	troca, err := h.trocaService.ObterPorVendaID(
		c.Request.Context(),
		vendaID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if troca == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Nenhuma troca vinculada a esta venda",
		})
		return
	}

	c.JSON(http.StatusOK, troca)
}

// GET /api/v1/trocas?limit=10&offset=0
func (h *TrocaController) Listar(c *gin.Context) {
	limit, offset := parsePaginacao(c)

	trocas, err := h.trocaService.Listar(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trocas)
}

// GET /api/v1/trocas/ultimo-mes?limit=10&offset=0
func (h *TrocaController) ListarUltimoMes(c *gin.Context) {
	limit, offset := parsePaginacao(c)

	trocas, err := h.trocaService.ListarUltimoMes(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trocas)
}

// GET /api/v1/trocas/tres-meses?limit=10&offset=0
func (h *TrocaController) ListarUltimosTresMeses(c *gin.Context) {
	limit, offset := parsePaginacao(c)

	trocas, err := h.trocaService.ListarUltimosTresMeses(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, trocas)
}

// GET /api/v1/trocas/seis-meses?limit=10&offset=0
func (h *TrocaController) ListarUltimosSeisMeses(c *gin.Context) {
	limit, offset := parsePaginacao(c)

	trocas, err := h.trocaService.ListarUltimosSeisMeses(c.Request.Context(), limit, offset)
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

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

package controllers

import (
	"context"
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

// POST /trocas
func (h *TrocaController) Criar(c *gin.Context) {
	empresaID, err := h.obterEmpresaID(c)
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

	// Injeta o ID da empresa recuperado pelo getter
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
	empresaID, err := h.obterEmpresaID(c)
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
	empresaID, err := h.obterEmpresaID(c)
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

// ==========================================
// GETTERS E FUNÇÕES AUXILIARES (HELPERS)
// ==========================================

// Função para recuperar o ID da empresa (Getter)
func (h *TrocaController) obterEmpresaID(c *gin.Context) (int64, error) {
	// OPÇÃO 1: Pegar de um contexto definido por um middleware de autenticação JWT
	// Se você tiver um middleware fazendo `c.Set("empresaID", id)`:
	if idVal, exists := c.Get("empresaID"); exists {
		// Tenta fazer o assert dependendo do tipo salvo (pode ser int, int64, float64 se vier de JSON)
		switch v := idVal.(type) {
		case int64:
			return v, nil
		case int:
			return int64(v), nil
		case float64:
			return int64(v), nil
		}
	}

	// OPÇÃO 2: Pegar direto de um Header HTTP (Ex: X-Empresa-ID)
	// Muito útil caso o tenant seja passado pelo front-end nos headers
	headerVal := c.GetHeader("X-Empresa-ID")
	if headerVal != "" {
		id, err := strconv.ParseInt(headerVal, 10, 64)
		if err == nil {
			return id, nil
		}
	}

	return 0, http.ErrNoCookie // Ou crie um erro customizado, ex: errors.New("empresaID não encontrado")
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

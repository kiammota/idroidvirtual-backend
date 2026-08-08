package controllers

import (
	"net/http"
	"strconv"
	"time"

	"backend/gateways"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type DespesaController struct {
	service *services.DespesaService
}

func NewDespesaController(service *services.DespesaService) *DespesaController {
	return &DespesaController{service: service}
}

// Criar lida com POST /despesas
func (c *DespesaController) Criar(ctx *gin.Context) {
	var despesa models.Despesa
	if err := ctx.ShouldBindJSON(&despesa); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Payload inválido: "+err.Error())
		return
	}

	empresaID := ctx.GetInt64("empresa_id")
	if empresaID != 0 {
		despesa.EmpresaID = empresaID
	}

	if err := c.service.Criar(ctx.Request.Context(), &despesa); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, despesa, "Despesa criada com sucesso")
}

// BuscarPorID lida com GET /despesas/:id
func (c *DespesaController) BuscarPorID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "ID inválido")
		return
	}

	empresaID := ctx.GetInt64("empresa_id")
	despesa, err := c.service.BuscarPorID(ctx.Request.Context(), id, empresaID)
	if err != nil {
		gateways.Error(ctx, http.StatusNotFound, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, despesa, "Despesa encontrada com sucesso")
}

// Listar lida com GET /despesas?inicio=YYYY-MM-DD&fim=YYYY-MM-DD
func (c *DespesaController) Listar(ctx *gin.Context) {
	empresaID := ctx.GetInt64("empresa_id")
	inicioStr := ctx.Query("inicio")
	fimStr := ctx.Query("fim")

	var inicio, fim time.Time
	var err error

	if inicioStr != "" {
		inicio, err = time.Parse("2006-01-02", inicioStr)
		if err != nil {
			gateways.Error(ctx, http.StatusBadRequest, "Formato de data de início inválido (use YYYY-MM-DD)")
			return
		}
	} else {
		agora := time.Now()
		inicio = time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	}

	if fimStr != "" {
		fimRaw, err := time.Parse("2006-01-02", fimStr)
		if err != nil {
			gateways.Error(ctx, http.StatusBadRequest, "Formato de data final inválido (use YYYY-MM-DD)")
			return
		}
		fim = time.Date(fimRaw.Year(), fimRaw.Month(), fimRaw.Day(), 23, 59, 59, 999999999, fimRaw.Location())
	} else {
		agora := time.Now()
		fim = time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	}

	despesas, err := c.service.ListarPorPeriodo(ctx.Request.Context(), empresaID, inicio, fim)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, despesas, "Despesas listadas com sucesso")
}

// Atualizar lida com PUT /despesas/:id
func (c *DespesaController) Atualizar(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "ID inválido")
		return
	}

	var despesa models.Despesa
	if err := ctx.ShouldBindJSON(&despesa); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Payload inválido: "+err.Error())
		return
	}

	despesa.ID = id
	despesa.EmpresaID = ctx.GetInt64("empresa_id")

	if err := c.service.Atualizar(ctx.Request.Context(), &despesa); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Despesa atualizada com sucesso")
}

// Deletar lida com DELETE /despesas/:id
func (c *DespesaController) Deletar(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "ID inválido")
		return
	}

	empresaID := ctx.GetInt64("empresa_id")
	if err := c.service.Deletar(ctx.Request.Context(), id, empresaID); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Despesa removida com sucesso")
}

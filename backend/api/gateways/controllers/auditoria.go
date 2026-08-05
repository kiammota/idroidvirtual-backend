package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type AuditoriaController struct {
	service *services.AuditoriaService
}

func NewAuditoriaController(service *services.AuditoriaService) *AuditoriaController {
	return &AuditoriaController{service: service}
}

func (c *AuditoriaController) ListarAuditorias(ctx *gin.Context) {
	periodoParam := ctx.DefaultQuery("periodo", "dia")
	paginaParam := ctx.DefaultQuery("pagina", "1")
	itensPorPaginaParam := ctx.DefaultQuery("itens_por_pagina", "20")

	pagina, err := strconv.Atoi(paginaParam)
	if err != nil || pagina < 1 {
		pagina = 1
	}

	itensPorPagina, err := strconv.Atoi(itensPorPaginaParam)
	if err != nil || itensPorPagina <= 0 {
		itensPorPagina = 20
	}

	registros, err := c.service.ListarPorPeriodo(
		ctx.Request.Context(),
		services.PeriodoTipo(periodoParam),
		pagina,
		itensPorPagina,
	)

	if err != nil {
		gateways.Error(
			ctx,
			http.StatusInternalServerError,
			"Erro ao buscar registros de auditoria: "+err.Error(),
		)
		return
	}

	ctx.JSON(http.StatusOK, registros)
}

package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"backend/gateways"
	"backend/models"
	"backend/repositories"

	"github.com/gin-gonic/gin"
)

type ProductController struct {
	repo *repositories.ProdutoRepository
}

func NewProductController(repo *repositories.ProdutoRepository) *ProductController {
	return &ProductController{repo: repo}
}

// GET /products
func (c *ProductController) ListAll(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	produtos, err := c.repo.Listar(ctx.Request.Context(), limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produtos, "Produtos listados com sucesso")
}

// POST /products
func (c *ProductController) Create(ctx *gin.Context) {
	var produto models.Produto

	if err := ctx.ShouldBindJSON(&produto); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos para criação do produto")
		return
	}

	if err := c.repo.Criar(ctx.Request.Context(), &produto); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, produto, "Produto criado com sucesso")
}

// GET /products/category/:category
func (c *ProductController) GetByCategory(ctx *gin.Context) {
	categoria := strings.TrimSpace(ctx.Param("category"))
	if categoria == "" {
		gateways.Error(ctx, http.StatusBadRequest, "Categoria é obrigatória")
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	produtos, err := c.repo.BuscarPorCategoria(ctx.Request.Context(), categoria, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produtos, "Produtos recuperados com sucesso")
}

// POST /products/:id (Atualização de produto)
func (c *ProductController) Update(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID do produto é inválido")
		return
	}

	var produto models.Produto
	if err := ctx.ShouldBindJSON(&produto); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos para atualização do produto")
		return
	}

	produto.ID = id

	if err := c.repo.Atualizar(ctx.Request.Context(), &produto); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produto, "Produto atualizado com sucesso")
}
func (c *ProductController) SetInventory(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID requisitado é inválido!")
		return
	}

	idQuant := ctx.Param("num")
	quant, err := strconv.ParseInt(idQuant, 10, 64)
	if err != nil || quant < 0 {
		gateways.Error(ctx, http.StatusBadRequest, "Quantidade é inválida!")
		return
	}

	err = c.repo.AlterarEstoque(ctx, id, int(quant))
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, 200, nil, "Estoque alterado com sucesso!")
}

// DELETE /products/:id
func (c *ProductController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID do produto é inválido")
		return
	}

	if err := c.repo.Deletar(ctx.Request.Context(), id); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Produto removido com sucesso")
}

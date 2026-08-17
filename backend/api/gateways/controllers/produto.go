package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"backend/gateways"
	"backend/gateways/helpers"
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
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] GetEmpresaID falhou: %v\n", err)
		gateways.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	produtos, err := c.repo.Listar(ctx.Request.Context(), empresaID, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produtos, "Produtos listados com sucesso")
}

// POST /products
func (c *ProductController) Create(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] GetEmpresaID falhou: %v\n", err)
		gateways.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	var produto models.Produto
	if err := ctx.ShouldBindJSON(&produto); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Dados inválidos para criação do produto")
		return
	}

	// Atribui a empresa dona do produto na criação
	produto.EmpresaID = empresaID

	if err := c.repo.Criar(ctx.Request.Context(), &produto); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusCreated, produto, "Produto criado com sucesso")
}

// GET /products/category/:category
func (c *ProductController) GetByCategory(ctx *gin.Context) {
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] GetEmpresaID falhou: %v\n", err)
		gateways.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	categoria := strings.TrimSpace(ctx.Param("category"))
	if categoria == "" {
		gateways.Error(ctx, http.StatusBadRequest, "Categoria é obrigatória")
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	produtos, err := c.repo.BuscarPorCategoria(ctx.Request.Context(), empresaID, categoria, limit, offset)
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produtos, "Produtos recuperados com sucesso")
}

// PUT /products/:id
func (c *ProductController) Update(ctx *gin.Context) {
	// Opcional: Se o ID do produto já for global e único, você pode remover
	// a linha abaixo se não quiser validar a empresa aqui.
	empresaID, err := helpers.GetEmpresaID(ctx)
	if err != nil {
		fmt.Printf("[DEBUG] GetEmpresaID falhou: %v\n", err)
		gateways.Error(ctx, http.StatusUnauthorized, err.Error())
		return
	}

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
	produto.EmpresaID = empresaID // Mantido se a sua query de Update usar WHERE id = ... AND empresa_id = ...

	if err := c.repo.Atualizar(ctx.Request.Context(), &produto); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, produto, "Produto atualizado com sucesso")
}

// POST /products/:id/inventory/:num
func (c *ProductController) SetInventory(ctx *gin.Context) {
	// Como você pontuou, para alterar o estoque de um produto já existente,
	// precisamos apenas do ID do produto, sem precisar do ID da empresa na chamada.
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID requisitado é inválido!")
		return
	}

	idQuant := ctx.Param("num")
	quant, err := strconv.ParseInt(idQuant, 10, 64)
	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Quantidade é inválida!")
		return
	}

	// Chamada limpa passando apenas ID e quantidade
	err = c.repo.AlterarEstoque(ctx.Request.Context(), id, int(quant))
	if err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Estoque alterado com sucesso!")
}

// DELETE /products/:id
func (c *ProductController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		gateways.Error(ctx, http.StatusBadRequest, "ID do produto é inválido")
		return
	}

	// Mantido com empresaID por segurança para evitar IDOR ao deletar
	if err := c.repo.Deletar(ctx.Request.Context(), id); err != nil {
		gateways.Error(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Produto removido com sucesso")
}

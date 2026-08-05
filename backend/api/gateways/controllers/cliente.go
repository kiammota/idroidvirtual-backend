package controllers

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"backend/models"

	"github.com/gin-gonic/gin"
)

// Interface ajustada conforme os métodos exatos do seu ClienteRepository
type ClientRepository interface {
	Criar(ctx context.Context, c *models.Cliente) error
	BuscarPorID(ctx context.Context, id int64) (*models.Cliente, error)
	Listar(ctx context.Context, limit, offset int) ([]models.Cliente, error)
	BuscarPorNome(ctx context.Context, nome string, limit, offset int) ([]models.Cliente, error)
	BuscarPorTelefone(ctx context.Context, telefone string) ([]models.Cliente, error)
	BuscarPorEndereco(ctx context.Context, endereco string) ([]models.Cliente, error)
	BuscarPorEmail(ctx context.Context, email string) ([]models.Cliente, error)
}

type ClientController struct {
	repo ClientRepository
}

func NewClientController(repo ClientRepository) *ClientController {
	return &ClientController{repo: repo}
}

// POST /clients/
func (c *ClientController) Create(ctx *gin.Context) {
	var client models.Cliente

	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Dados inválidos para criação do cliente",
		})
		return
	}

	if err := c.repo.Criar(ctx.Request.Context(), &client); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    client,
		"message": "Cliente criado com sucesso",
	})
}

// GET /clients/
func (c *ClientController) ListAll(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	clientes, err := c.repo.Listar(ctx.Request.Context(), limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientes,
	})
}

// GET /clients/:id
func (c *ClientController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID do cliente é inválido",
		})
		return
	}

	cliente, err := c.repo.BuscarPorID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	if cliente == nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Cliente não encontrado",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    cliente,
	})
}

// GET /clients/number?telefone=...
func (c *ClientController) SearchByPhone(ctx *gin.Context) {
	phone := strings.TrimSpace(ctx.Query("telefone"))
	if phone == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parâmetro 'telefone' é obrigatório",
		})
		return
	}

	clientes, err := c.repo.BuscarPorTelefone(ctx.Request.Context(), phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientes,
	})
}

// GET /clients/address?endereco=...
func (c *ClientController) SearchByAddress(ctx *gin.Context) {
	address := strings.TrimSpace(ctx.Query("endereco"))
	if address == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parâmetro 'endereco' é obrigatório",
		})
		return
	}

	clientes, err := c.repo.BuscarPorEndereco(ctx.Request.Context(), address)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientes,
	})
}

// GET /clients/name?nome=...
func (c *ClientController) SearchByName(ctx *gin.Context) {
	nome := strings.TrimSpace(ctx.Query("nome"))
	if nome == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parâmetro 'nome' é obrigatório",
		})
		return
	}

	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))

	clientes, err := c.repo.BuscarPorNome(ctx.Request.Context(), nome, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientes,
	})
}

// GET /clients/email?email=...
func (c *ClientController) SearchByEmail(ctx *gin.Context) {
	email := strings.TrimSpace(ctx.Query("email"))
	if email == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Parâmetro 'email' é obrigatório",
		})
		return
	}

	clientes, err := c.repo.BuscarPorEmail(ctx.Request.Context(), email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    clientes,
	})
}

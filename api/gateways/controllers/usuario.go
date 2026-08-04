package controllers

import (
	"net/http"
	"strconv"

	"backend/gateways"
	"backend/models"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type UsuarioController struct {
	UsuarioService *services.UsuarioService
	LoginService   *services.LoginService
}

func NewUsuarioController(
	usuarioService *services.UsuarioService,
	loginService *services.LoginService,
) *UsuarioController {

	return &UsuarioController{
		UsuarioService: usuarioService,
		LoginService:   loginService,
	}
}

// POST /usuarios
func (c *UsuarioController) CriarUsuario(ctx *gin.Context) {

	var usuario models.Usuario

	if err := ctx.ShouldBindJSON(&usuario); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "dados inválidos")
		return
	}

	err := c.UsuarioService.Criar(&usuario)

	if err != nil {
		gateways.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	gateways.Success(
		ctx,
		http.StatusCreated,
		usuario,
		"usuário criado com sucesso",
	)
}

func (c *UsuarioController) ObterMeuUsuario(ctx *gin.Context) {

	usuarioID, existe := ctx.Get("usuario_id")

	if !existe {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "usuário não autenticado",
		})
		return
	}

	usuario, err := c.UsuarioService.BuscarPorID(
		usuarioID.(int64),
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, usuario)
}

// GET /usuarios/:id
func (c *UsuarioController) BuscarUsuario(ctx *gin.Context) {

	id, err := strconv.ParseInt(
		ctx.Param("id"),
		10,
		64,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "id inválido",
		})
		return
	}

	usuario, err := c.UsuarioService.BuscarPorID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, usuario)
}

func (c *UsuarioController) UsuarioSenha(ctx *gin.Context) {
	var body struct {
		Senha string `json:"senha" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&body); err != nil {
		gateways.Error(ctx, http.StatusBadRequest, "Senha obrigatória")
		return
	}

	if body.Senha != "1123" {
		gateways.Error(ctx, http.StatusUnauthorized, "Senha administrativa inválida")
		return
	}

	gateways.Success(ctx, http.StatusOK, nil, "Senha administrativa validada")
}

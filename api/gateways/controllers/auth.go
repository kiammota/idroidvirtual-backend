package controllers

import (
	"net/http"

	"backend/models"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	LoginService  *services.LoginService
	AuditoriaRepo *repositories.AuditoriaRepository
}

func NewAuthController(
	loginService *services.LoginService, auditoria *repositories.AuditoriaRepository,
) *AuthController {

	return &AuthController{
		LoginService:  loginService,
		AuditoriaRepo: auditoria,
	}
}

func (c *AuthController) UsuarioLogar(ctx *gin.Context) {

	var login models.LoginRequest

	if err := ctx.ShouldBindJSON(&login); err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"message": "Email e senha obrigatórios",
			},
		)
		return
	}

	response, err := c.LoginService.Login(
		login.Email,
		login.Senha,
	)

	if err != nil {
		ctx.AbortWithStatusJSON(
			http.StatusUnauthorized,
			gin.H{
				"message": err.Error(),
			},
		)
		return
	}

	err = c.AuditoriaRepo.Criar(
		ctx,
		models.Auditoria{
			UsuarioID:  response.Usuario.ID,
			Acao:       "LOGIN",
			Tabela:     "usuarios",
			RegistroID: response.Usuario.ID,

			ValoresNovos: map[string]interface{}{
				"email": response.Usuario.Email,
			},
		},
	)

	if err != nil {
		ctx.Error(err)
	}

	ctx.JSON(
		http.StatusOK,
		response,
	)
}

func (c *AuthController) UsuarioSenha(ctx *gin.Context) {
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

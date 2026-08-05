package controllers

import (
	"net/http"

	"backend/gateways"
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
	loginService *services.LoginService,
	auditoria *repositories.AuditoriaRepository,
) *AuthController {

	return &AuthController{
		LoginService:  loginService,
		AuditoriaRepo: auditoria,
	}
}

func (c *AuthController) UsuarioLogar(ctx *gin.Context) {

	var login models.LoginRequest

	if err := ctx.ShouldBindJSON(&login); err != nil {
		gateways.Error(
			ctx,
			http.StatusBadRequest,
			"email e senha obrigatórios",
		)
		return
	}

	response, err := c.LoginService.Login(
		login.Email,
		login.Senha,
	)

	if err != nil {
		gateways.Error(
			ctx,
			http.StatusUnauthorized,
			err.Error(),
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
		gateways.Error(
			ctx,
			http.StatusInternalServerError,
			"erro ao registrar auditoria",
		)
		return
	}

	gateways.Success(
		ctx,
		http.StatusOK,
		response.Token,
		"login realizado com sucesso",
	)
}

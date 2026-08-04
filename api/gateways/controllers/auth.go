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

package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	UsuarioRepository *repositories.UsuarioRepository
	TokenService      *TokenService
}

func NewLoginService(
	usuarioRepository *repositories.UsuarioRepository,
	tokenService *TokenService,
) *LoginService {

	return &LoginService{
		UsuarioRepository: usuarioRepository,
		TokenService:      tokenService,
	}
}

type LoginResponse struct {
	Token   string          `json:"token"`
	Usuario *models.Usuario `json:"usuario"`
}

// Realiza login
func (s *LoginService) Login(email string, senha string) (*LoginResponse, error) {

	if email == "" {
		return nil, errors.New("email obrigatório")
	}

	if senha == "" {
		return nil, errors.New("senha obrigatória")
	}

	usuario, err := s.UsuarioRepository.BuscarPorEmail(email)

	if err != nil {
		return nil, errors.New("email ou senha inválidos")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(usuario.SenhaHash),
		[]byte(senha),
	)

	if err != nil {
		return nil, errors.New("email ou senha inválidos")
	}

	token, err := s.TokenService.GerarToken(
		usuario.ID,
		usuario.Email,
	)

	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token:   token,
		Usuario: usuario,
	}, nil
}

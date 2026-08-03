package services

import (
	"backend/models"
	"backend/repositories"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UsuarioService struct {
	UsuarioRepository *repositories.UsuarioRepository
}

func NewUsuarioService(
	usuarioRepository *repositories.UsuarioRepository,
) *UsuarioService {

	return &UsuarioService{
		UsuarioRepository: usuarioRepository,
	}
}

// Busca usuário pelo ID
func (s *UsuarioService) BuscarPorID(id int64) (*models.Usuario, error) {

	usuario, err := s.UsuarioRepository.BuscarPorID(id)

	if err != nil {
		return nil, err
	}

	return usuario, nil
}

// Busca usuário pelo email
func (s *UsuarioService) BuscarPorEmail(email string) (*models.Usuario, error) {

	if email == "" {
		return nil, errors.New("email obrigatório")
	}

	return s.UsuarioRepository.BuscarPorEmail(email)
}

// Cria um usuário
func (s *UsuarioService) Criar(usuario *models.Usuario) error {

	if usuario.Nome == "" {
		return errors.New("nome obrigatório")
	}

	if usuario.Email == "" {
		return errors.New("email obrigatório")
	}

	if usuario.SenhaHash == "" {
		return errors.New("senha obrigatória")
	}

	senhaHash, err := bcrypt.GenerateFromPassword(
		[]byte(usuario.SenhaHash),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	usuario.SenhaHash = string(senhaHash)

	return s.UsuarioRepository.Criar(usuario)
}

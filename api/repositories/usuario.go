package repositories

import (
	"backend/models"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsuarioRepository struct {
	DB *pgxpool.Pool
}

func NewUsuarioRepository(db *pgxpool.Pool) *UsuarioRepository {
	return &UsuarioRepository{
		DB: db,
	}
}

func (r *UsuarioRepository) BuscarPorEmail(email string) (*models.Usuario, error) {

	query := `
		SELECT 
			id,
			empresa_id,
			email,
			nome,
			senha,
			criado_em
		FROM usuarios
		WHERE email = $1
		LIMIT 1
	`

	usuario := &models.Usuario{}

	err := r.DB.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(
		&usuario.ID,
		&usuario.EmpresaID,
		&usuario.Email,
		&usuario.Nome,
		&usuario.SenhaHash,
		&usuario.CriadoEm,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuário não encontrado")
		}

		return nil, err
	}

	return usuario, nil
}

func (r *UsuarioRepository) BuscarPorID(id int64) (*models.Usuario, error) {

	query := `
		SELECT 
			id,
			empresa_id,
			email,
			nome,
			senha,
			criado_em
		FROM usuarios
		WHERE id = $1
		LIMIT 1
	`

	usuario := &models.Usuario{}

	err := r.DB.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&usuario.ID,
		&usuario.EmpresaID,
		&usuario.Email,
		&usuario.Nome,
		&usuario.SenhaHash,
		&usuario.CriadoEm,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usuário não encontrado")
		}

		return nil, err
	}

	return usuario, nil
}

func (r *UsuarioRepository) Criar(usuario *models.Usuario) error {

	query := `
		INSERT INTO usuarios (
			empresa_id,
			email,
			nome,
			senha,
			criado_em
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			NOW()
		)
		RETURNING id, criado_em
	`

	return r.DB.QueryRow(
		context.Background(),
		query,
		usuario.EmpresaID,
		usuario.Email,
		usuario.Nome,
		usuario.SenhaHash,
	).Scan(
		&usuario.ID,
		&usuario.CriadoEm,
	)
}

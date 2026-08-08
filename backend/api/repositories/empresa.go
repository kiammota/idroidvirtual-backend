package repositories

import (
	"context"
	"fmt"

	"backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmpresaRepository struct {
	Pool *pgxpool.Pool
}

func NewEmpresaRepository(pool *pgxpool.Pool) *EmpresaRepository {
	return &EmpresaRepository{Pool: pool}
}

// BuscarPorUsuarioID retorna todos os dados da empresa vinculada ao ID do usuário
func (r *EmpresaRepository) BuscarPorUsuarioID(ctx context.Context, usuarioID int64) (*models.Empresa, error) {
	query := `
		SELECT 
			e.id, 
			e.nome, 
			e.cnpj, 
			e.criado_em 
		FROM empresas e
		INNER JOIN usuarios u ON u.empresa_id = e.id
		WHERE u.id = $1
	`

	var empresa models.Empresa

	err := r.Pool.QueryRow(ctx, query, usuarioID).Scan(
		&empresa.ID,
		&empresa.Nome,
		&empresa.CNPJ,
		&empresa.CriadoEm,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar empresa do usuário #%d: %w", usuarioID, err)
	}

	return &empresa, nil
}

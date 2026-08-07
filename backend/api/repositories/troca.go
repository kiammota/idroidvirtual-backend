package repositories

import (
	"context"
	"errors"

	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TrocaRepository struct {
	Pool *pgxpool.Pool
}

func NewTrocaRepository(pool *pgxpool.Pool) TrocaRepository {
	return TrocaRepository{
		Pool: pool,
	}
}

// Criar insere um registro de troca
func (r *TrocaRepository) Criar(
	ctx context.Context,
	t *models.Troca,
) error {

	query := `
		INSERT INTO trocas (
			venda_id,
			modelo_aparelho,
			imei,
			valor_avaliacao,
			condicao,
			observacoes
		)
		VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, criado_em
	`

	if t.Condicao == "" {
		t.Condicao = "USADO"
	}

	return r.Pool.QueryRow(
		ctx,
		query,
		t.VendaID,
		t.ModeloAparelho,
		t.IMEI,
		t.ValorAvaliacao,
		t.Condicao,
		t.Observacoes,
	).Scan(
		&t.ID,
		&t.CriadoEm,
	)
}

// BuscarPorVendaID retorna a troca vinculada a uma venda
func (r *TrocaRepository) BuscarPorVendaID(
	ctx context.Context,
	vendaID int64,
) (*models.Troca, error) {

	query := `
		SELECT
			id,
			venda_id,
			modelo_aparelho,
			imei,
			valor_avaliacao,
			condicao,
			observacoes,
			criado_em
		FROM trocas
		WHERE venda_id = $1
	`

	var t models.Troca

	err := r.Pool.QueryRow(
		ctx,
		query,
		vendaID,
	).Scan(
		&t.ID,
		&t.VendaID,
		&t.ModeloAparelho,
		&t.IMEI,
		&t.ValorAvaliacao,
		&t.Condicao,
		&t.Observacoes,
		&t.CriadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &t, nil
}

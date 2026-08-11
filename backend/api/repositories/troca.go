package repositories

import (
	"context"
	"errors"
	"time"

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

func (r *TrocaRepository) Listar(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {

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
		ORDER BY criado_em DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trocas []*models.Troca

	for rows.Next() {
		var t models.Troca
		err := rows.Scan(
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
			return nil, err
		}
		trocas = append(trocas, &t)
	}

	return trocas, rows.Err()
}

// ListarPorIntervaloMeses busca trocas criadas nos últimos N meses com paginação
func (r *TrocaRepository) ListarPorIntervaloMeses(
	ctx context.Context,
	meses int,
	limit, offset int,
) ([]*models.Troca, error) {

	dataInicio := time.Now().AddDate(0, -meses, 0)

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
		WHERE criado_em >= $1
		ORDER BY criado_em DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, dataInicio, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trocas []*models.Troca

	for rows.Next() {
		var t models.Troca
		err := rows.Scan(
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
			return nil, err
		}
		trocas = append(trocas, &t)
	}

	return trocas, rows.Err()
}

// ListarUltimoMes retorna as trocas do último mês
func (r *TrocaRepository) ListarUltimoMes(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, 1, limit, offset)
}

// ListarUltimosTresMeses retorna as trocas dos últimos 3 meses
func (r *TrocaRepository) ListarUltimosTresMeses(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, 3, limit, offset)
}

// ListarUltimosSeisMeses retorna as trocas dos últimos 6 meses
func (r *TrocaRepository) ListarUltimosSeisMeses(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, 6, limit, offset)
}

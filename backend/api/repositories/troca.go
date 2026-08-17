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
			empresa_id,
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
		t.EmpresaID,
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

// BuscarPorID retorna uma troca pelo seu ID e ID da Empresa
func (r *TrocaRepository) BuscarPorID(
	ctx context.Context,
	id int64,
	empresaID int64,
) (*models.Troca, error) {

	query := `
		SELECT
			id,
			empresa_id,
			modelo_aparelho,
			imei,
			valor_avaliacao,
			condicao,
			observacoes,
			criado_em
		FROM trocas
		WHERE id = $1 AND empresa_id = $2
	`

	var t models.Troca

	err := r.Pool.QueryRow(
		ctx,
		query,
		id,
		empresaID,
	).Scan(
		&t.ID,
		&t.EmpresaID,
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

// Listar retorna uma lista paginada de trocas da empresa
func (r *TrocaRepository) Listar(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {

	query := `
		SELECT
			id,
			empresa_id,
			modelo_aparelho,
			imei,
			valor_avaliacao,
			condicao,
			observacoes,
			criado_em
		FROM trocas
		WHERE empresa_id = $1
		ORDER BY criado_em DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trocas []*models.Troca

	for rows.Next() {
		var t models.Troca
		err := rows.Scan(
			&t.ID,
			&t.EmpresaID,
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

// ListarPorIntervaloMeses busca trocas de uma empresa criadas nos últimos N meses com paginação
func (r *TrocaRepository) ListarPorIntervaloMeses(
	ctx context.Context,
	empresaID int64,
	meses int,
	limit, offset int,
) ([]*models.Troca, error) {

	dataInicio := time.Now().AddDate(0, -meses, 0)

	query := `
		SELECT
			id,
			empresa_id,
			modelo_aparelho,
			imei,
			valor_avaliacao,
			condicao,
			observacoes,
			criado_em
		FROM trocas
		WHERE empresa_id = $1 AND criado_em >= $2
		ORDER BY criado_em DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, dataInicio, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trocas []*models.Troca

	for rows.Next() {
		var t models.Troca
		err := rows.Scan(
			&t.ID,
			&t.EmpresaID,
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

// ListarUltimoMes retorna as trocas da empresa do último mês
func (r *TrocaRepository) ListarUltimoMes(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, empresaID, 1, limit, offset)
}

// ListarUltimosTresMeses retorna as trocas da empresa dos últimos 3 meses
func (r *TrocaRepository) ListarUltimosTresMeses(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, empresaID, 3, limit, offset)
}

// ListarUltimosSeisMeses retorna as trocas da empresa dos últimos 6 meses
func (r *TrocaRepository) ListarUltimosSeisMeses(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {
	return r.ListarPorIntervaloMeses(ctx, empresaID, 6, limit, offset)
}

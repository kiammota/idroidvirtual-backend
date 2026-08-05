package repositories

import (
	"context"
	"errors"
	"time"

	"backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DespesaRepository struct {
	Pool *pgxpool.Pool
}

func NewDespesaRepository(pool *pgxpool.Pool) DespesaRepository {
	return DespesaRepository{Pool: pool}
}

// Criar regista uma nova despesa
func (r *DespesaRepository) Criar(ctx context.Context, d *models.Despesa) error {
	query := `
		INSERT INTO despesas (
			usuario_id, descricao, valor, categoria, status, data_vencimento, data_pagamento
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, criado_em
	`

	return r.Pool.QueryRow(
		ctx,
		query,
		d.UsuarioID,
		d.Descricao,
		d.Valor,
		d.Categoria,
		d.Status,
		d.DataVencimento,
		d.DataPagamento,
	).Scan(&d.ID, &d.CriadoEm)
}

// MarcarComoPaga atualiza a despesa para PAGO e regista a data de pagamento
func (r *DespesaRepository) MarcarComoPaga(ctx context.Context, id int64) error {
	query := `
		UPDATE despesas
		SET status = 'PAGO',
		    data_pagamento = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhuma despesa encontrada para atualizar")
	}

	return nil
}

// ListarPendentes busca despesas que ainda não foram pagas (útil para lembretes de aluguel/contas)
func (r *DespesaRepository) ListarPendentes(ctx context.Context) ([]models.Despesa, error) {
	query := `
		SELECT id, usuario_id, descricao, valor, categoria, status, data_vencimento, data_pagamento, criado_em
		FROM despesas
		WHERE status = 'PENDENTE'
		ORDER BY data_vencimento ASC
	`

	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var despesas []models.Despesa
	for rows.Next() {
		var d models.Despesa
		err := rows.Scan(
			&d.ID,
			&d.UsuarioID,
			&d.Descricao,
			&d.Valor,
			&d.Categoria,
			&d.Status,
			&d.DataVencimento,
			&d.DataPagamento,
			&d.CriadoEm,
		)
		if err != nil {
			return nil, err
		}
		despesas = append(despesas, d)
	}

	return despesas, nil
}

// ObterTotalDespesasPagasPorPeriodo soma o total pago entre duas datas (usado no Fechamento Financeiro)
func (r *DespesaRepository) ObterTotalDespesasPagasPorPeriodo(ctx context.Context, inicio, fim time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(valor), 0)
		FROM despesas
		WHERE status = 'PAGO' AND data_pagamento >= $1 AND data_pagamento <= $2
	`

	var total float64
	err := r.Pool.QueryRow(ctx, query, inicio, fim).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

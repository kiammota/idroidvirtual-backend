package repositories

import (
	"backend/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VendaRepository struct {
	Pool *pgxpool.Pool
}

func NewVendaRepository(pool *pgxpool.Pool) VendaRepository {
	return VendaRepository{Pool: pool}
}

// CriarTx insere a venda com o status inicial 'CONCLUIDA' usando transação
func (r *VendaRepository) CriarTx(ctx context.Context, tx pgx.Tx, v *models.Venda) error {
	query := `
		INSERT INTO vendas (
			empresa_id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, criado_em
	`

	if v.Status == "" {
		v.Status = "CONCLUIDA"
	}

	return tx.QueryRow(
		ctx,
		query,
		v.EmpresaID,
		v.ProdutoID,
		v.UsuarioID,
		v.Quantidade,
		v.Total,
		v.TipoPagamento,
		v.Status,
	).Scan(&v.ID, &v.CriadoEm)
}

// Criar insere uma nova venda diretamente pelo pool
func (r *VendaRepository) Criar(ctx context.Context, v *models.Venda) error {
	query := `
		INSERT INTO vendas (
			empresa_id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, criado_em
	`

	if v.Status == "" {
		v.Status = "CONCLUIDA"
	}

	return r.Pool.QueryRow(
		ctx,
		query,
		v.EmpresaID,
		v.ProdutoID,
		v.UsuarioID,
		v.Quantidade,
		v.Total,
		v.TipoPagamento,
		v.Status,
	).Scan(&v.ID, &v.CriadoEm)
}

// BuscarPorID retorna uma venda específica pelo ID e pela Empresa
func (r *VendaRepository) BuscarPorID(ctx context.Context, id int64) (*models.Venda, error) {
	query := `
		SELECT id, empresa_id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE id = $1
	`

	var v models.Venda
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&v.ID,
		&v.EmpresaID,
		&v.ProdutoID,
		&v.UsuarioID,
		&v.Quantidade,
		&v.Total,
		&v.TipoPagamento,
		&v.Status,
		&v.CriadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &v, nil
}

// BuscarPorPeriodo retorna a lista de vendas de uma empresa realizadas entre duas datas
func (r *VendaRepository) BuscarPorPeriodo(ctx context.Context, empresaID int64, inicio, fim time.Time) ([]models.Venda, error) {
	query := `
		SELECT id, empresa_id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE empresa_id = $1 AND criado_em >= $2 AND criado_em <= $3
		ORDER BY criado_em DESC
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, inicio, fim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []models.Venda
	for rows.Next() {
		var v models.Venda
		err := rows.Scan(
			&v.ID,
			&v.EmpresaID,
			&v.ProdutoID,
			&v.UsuarioID,
			&v.Quantidade,
			&v.Total,
			&v.TipoPagamento,
			&v.Status,
			&v.CriadoEm,
		)
		if err != nil {
			return nil, err
		}
		vendas = append(vendas, v)
	}

	return vendas, nil
}

// Listar retorna o histórico de vendas de uma empresa com suporte a paginação
func (r *VendaRepository) Listar(ctx context.Context, empresaID int64, limit, offset int) ([]models.Venda, error) {
	query := `
		SELECT id, empresa_id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE empresa_id = $1
		ORDER BY criado_em DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []models.Venda
	for rows.Next() {
		var v models.Venda
		err := rows.Scan(
			&v.ID,
			&v.EmpresaID,
			&v.ProdutoID,
			&v.UsuarioID,
			&v.Quantidade,
			&v.Total,
			&v.TipoPagamento,
			&v.Status,
			&v.CriadoEm,
		)
		if err != nil {
			return nil, err
		}
		vendas = append(vendas, v)
	}

	return vendas, nil
}

// ObterFaturamentoMesAtual retorna a soma total das vendas no mês vigente para a empresa
func (r *VendaRepository) ObterFaturamentoMesAtual(ctx context.Context, empresaID int64) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total), 0)
		FROM vendas
		WHERE empresa_id = $1 AND criado_em >= DATE_TRUNC('month', CURRENT_DATE) AND status != 'CANCELADA'
	`

	var total float64
	err := r.Pool.QueryRow(ctx, query, empresaID).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// ObterFaturamentoPorPeriodo retorna o faturamento de uma empresa entre duas datas
func (r *VendaRepository) ObterFaturamentoPorPeriodo(ctx context.Context, empresaID int64, inicio, fim time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total), 0)
		FROM vendas
		WHERE empresa_id = $1 AND criado_em >= $2 AND criado_em <= $3 AND status != 'CANCELADA'
	`

	var total float64
	err := r.Pool.QueryRow(ctx, query, empresaID, inicio, fim).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// CancelarTx altera o status da venda para CANCELADA dentro da transação, garantindo que pertença à empresa
func (r *VendaRepository) CancelarTx(ctx context.Context, tx pgx.Tx, vendaID int64) error {
	query := `
		UPDATE vendas
		SET status = 'CANCELADA'
		WHERE id = $1 AND status != 'CANCELADA'
	`

	cmdTag, err := tx.Exec(ctx, query, vendaID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("venda não encontrada, já estava cancelada ou não pertence a esta empresa")
	}

	return nil
}

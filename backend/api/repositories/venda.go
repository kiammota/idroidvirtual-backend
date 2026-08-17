package repositories

import (
	"context"
	"errors"
	"time"

	"backend/models"

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
			produto_id, usuario_id, quantidade, total, tipo_pagamento, status
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, criado_em
	`

	if v.Status == "" {
		v.Status = "CONCLUIDA"
	}

	return tx.QueryRow(
		ctx,
		query,
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
			produto_id, usuario_id, quantidade, total, tipo_pagamento, status
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, criado_em
	`

	if v.Status == "" {
		v.Status = "CONCLUIDA"
	}

	return r.Pool.QueryRow(
		ctx,
		query,
		v.ProdutoID,
		v.UsuarioID,
		v.Quantidade,
		v.Total,
		v.TipoPagamento,
		v.Status,
	).Scan(&v.ID, &v.CriadoEm)
}

// BuscarPorID retorna uma venda específica pelo ID incluindo o status
func (r *VendaRepository) BuscarPorID(ctx context.Context, id int64) (*models.Venda, error) {
	query := `
		SELECT id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE id = $1
	`

	var v models.Venda
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&v.ID,
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

// BuscarPorPeriodo retorna a lista de vendas realizadas entre duas datas
func (r *VendaRepository) BuscarPorPeriodo(ctx context.Context, inicio, fim time.Time) ([]models.Venda, error) {
	query := `
		SELECT id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE criado_em >= $1 AND criado_em <= $2
		ORDER BY criado_em DESC
	`

	rows, err := r.Pool.Query(ctx, query, inicio, fim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []models.Venda
	for rows.Next() {
		var v models.Venda
		err := rows.Scan(
			&v.ID,
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

// Listar retorna o histórico de vendas com suporte a paginação
func (r *VendaRepository) Listar(ctx context.Context, limit, offset int) ([]models.Venda, error) {
	query := `
		SELECT id, produto_id, usuario_id, quantidade, total, tipo_pagamento, status, criado_em
		FROM vendas
		ORDER BY criado_em DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var vendas []models.Venda
	for rows.Next() {
		var v models.Venda
		err := rows.Scan(
			&v.ID,
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

// ObterFaturamentoMesAtual retorna a soma total das vendas no mês vigente (excluindo canceladas)
func (r *VendaRepository) ObterFaturamentoMesAtual(ctx context.Context) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total), 0)
		FROM vendas
		WHERE criado_em >= DATE_TRUNC('month', CURRENT_DATE) AND status != 'CANCELADA'
	`

	var total float64
	err := r.Pool.QueryRow(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// ObterFaturamentoPorPeriodo retorna o faturamento entre duas datas (excluindo canceladas)
func (r *VendaRepository) ObterFaturamentoPorPeriodo(ctx context.Context, inicio, fim time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total), 0)
		FROM vendas
		WHERE criado_em >= $1 AND criado_em <= $2 AND status != 'CANCELADA'
	`

	var total float64
	err := r.Pool.QueryRow(ctx, query, inicio, fim).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// CancelarTx altera o status da venda para CANCELADA dentro da transação
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
		return errors.New("venda não encontrada ou já estava cancelada")
	}

	return nil
}

// AumentarEstoqueTx incrementa o estoque dentro de uma transação apenas pelo ID do produto
func (r *ProdutoRepository) AumentarEstoqueTx(ctx context.Context, tx pgx.Tx, id int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser adicionada deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque + $1,
			atualizado_em = NOW()
		WHERE id = $2
	`

	cmdTag, err := tx.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para aumentar o estoque")
	}

	return nil
}

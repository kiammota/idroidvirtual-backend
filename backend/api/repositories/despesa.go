package repositories

import (
	"context"
	"fmt"
	"time"

	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DespesaRepository struct {
	pool *pgxpool.Pool
}

func NewDespesaRepository(pool *pgxpool.Pool) *DespesaRepository {
	return &DespesaRepository{
		pool: pool,
	}
}

// Criar insere uma nova despesa no banco de dados
func (r *DespesaRepository) Criar(ctx context.Context, d *models.Despesa) error {
	query := `
		INSERT INTO despesas (
			empresa_id, usuario_id, descricao, categoria, valor, valor_pago,
			forma_pagamento, status, data_vencimento, data_pagamento,
			comprovante_url, numero_documento
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, criado_em;
	`

	err := r.pool.QueryRow(
		ctx, query,
		d.EmpresaID, d.UsuarioID, d.Descricao, d.Categoria, d.Valor, d.ValorPago,
		d.FormaPagamento, d.Status, d.DataVencimento, d.DataPagamento,
		d.ComprovanteURL, d.NumeroDocumento,
	).Scan(&d.ID, &d.CriadoEm)

	if err != nil {
		return fmt.Errorf("erro ao criar despesa: %w", err)
	}

	return nil
}

// BuscarPorID retorna uma despesa pelo seu ID e ID da empresa
func (r *DespesaRepository) BuscarPorID(ctx context.Context, id, empresaID int64) (*models.Despesa, error) {
	query := `
		SELECT 
			id, empresa_id, usuario_id, descricao, categoria, valor, valor_pago,
			forma_pagamento, status, data_vencimento, data_pagamento,
			comprovante_url, numero_documento, criado_em
		FROM despesas
		WHERE id = $1 AND empresa_id = $2;
	`

	var d models.Despesa
	err := r.pool.QueryRow(ctx, query, id, empresaID).Scan(
		&d.ID, &d.EmpresaID, &d.UsuarioID, &d.Descricao, &d.Categoria, &d.Valor, &d.ValorPago,
		&d.FormaPagamento, &d.Status, &d.DataVencimento, &d.DataPagamento,
		&d.ComprovanteURL, &d.NumeroDocumento, &d.CriadoEm,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("despesa não encontrada")
		}
		return nil, fmt.Errorf("erro ao buscar despesa: %w", err)
	}

	return &d, nil
}

// ListarPorEmpresaEPeriodo busca despesas dentro de um intervalo de datas por vencimento
func (r *DespesaRepository) ListarPorEmpresaEPeriodo(ctx context.Context, empresaID int64, inicio, fim time.Time) ([]models.Despesa, error) {
	query := `
		SELECT 
			id, empresa_id, usuario_id, descricao, categoria, valor, valor_pago,
			forma_pagamento, status, data_vencimento, data_pagamento,
			comprovante_url, numero_documento, criado_em
		FROM despesas
		WHERE empresa_id = $1 
		  AND data_vencimento BETWEEN $2 AND $3
		ORDER BY data_vencimento DESC;
	`

	rows, err := r.pool.Query(ctx, query, empresaID, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar despesas por período: %w", err)
	}
	defer rows.Close()

	despesas := make([]models.Despesa, 0)
	for rows.Next() {
		var d models.Despesa
		err := rows.Scan(
			&d.ID, &d.EmpresaID, &d.UsuarioID, &d.Descricao, &d.Categoria, &d.Valor, &d.ValorPago,
			&d.FormaPagamento, &d.Status, &d.DataVencimento, &d.DataPagamento,
			&d.ComprovanteURL, &d.NumeroDocumento, &d.CriadoEm,
		)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler linha de despesa: %w", err)
		}
		despesas = append(despesas, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração de despesas: %w", err)
	}

	return despesas, nil
}

// Atualizar atualiza os dados de uma despesa existente
func (r *DespesaRepository) Atualizar(ctx context.Context, d *models.Despesa) error {
	query := `
		UPDATE despesas
		SET 
			usuario_id = $1,
			descricao = $2,
			categoria = $3,
			valor = $4,
			valor_pago = $5,
			forma_pagamento = $6,
			status = $7,
			data_vencimento = $8,
			data_pagamento = $9,
			comprovante_url = $10,
			numero_documento = $11
		WHERE id = $12 AND empresa_id = $13;
	`

	cmd, err := r.pool.Exec(
		ctx, query,
		d.UsuarioID, d.Descricao, d.Categoria, d.Valor, d.ValorPago,
		d.FormaPagamento, d.Status, d.DataVencimento, d.DataPagamento,
		d.ComprovanteURL, d.NumeroDocumento, d.ID, d.EmpresaID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar despesa: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("nenhuma despesa encontrada para atualizar")
	}

	return nil
}

// Deletar remove um registro de despesa pelo ID
func (r *DespesaRepository) Deletar(ctx context.Context, id, empresaID int64) error {
	query := `DELETE FROM despesas WHERE id = $1 AND empresa_id = $2;`

	cmd, err := r.pool.Exec(ctx, query, id, empresaID)
	if err != nil {
		return fmt.Errorf("erro ao deletar despesa: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("despesa não encontrada para exclusão")
	}

	return nil
}

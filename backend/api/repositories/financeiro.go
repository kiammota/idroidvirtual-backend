package repositories

import (
	"context"
	"fmt"
	"time"

	"backend/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FinanceiroRepository struct {
	pool *pgxpool.Pool
}

func NewFinanceiroRepository(pool *pgxpool.Pool) *FinanceiroRepository {
	return &FinanceiroRepository{
		pool: pool,
	}
}

// ObterResumoPorPeriodo executa uma única query otimizada agregando vendas, OS, trocas e despesas
func (r *FinanceiroRepository) ObterResumoPorPeriodo(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	query := `
		WITH 
		vendas_summary AS (
			SELECT 
				COALESCE(SUM(total), 0) AS total_vendas,
				COUNT(*) AS qtd_vendas
			FROM vendas
			WHERE status != 'CANCELADA'
			  AND criado_em BETWEEN $1 AND $2
		),
		os_summary AS (
			SELECT 
				COALESCE(SUM(valor), 0) AS total_servicos,
				COALESCE(SUM(custo), 0) AS custo_servicos,
				COUNT(*) AS qtd_os
			FROM ordens_servico
			WHERE status = 'FINALIZADA'
			  AND criado_em BETWEEN $1 AND $2
		),
		trocas_summary AS (
			SELECT 
				COALESCE(SUM(valor_avaliacao), 0) AS valor_estoque_perdido,
				COUNT(*) AS qtd_perdas
			FROM trocas
			WHERE criado_em BETWEEN $1 AND $2
		),
		despesas_summary AS (
			SELECT 
				COALESCE(SUM(valor_pago), 0) AS total_despesas
			FROM despesas
			WHERE status = 'PAGO'
			  AND data_pagamento BETWEEN $1 AND $2
		)
		SELECT 
			v.total_vendas,
			v.qtd_vendas,
			o.total_servicos,
			o.custo_servicos,
			o.qtd_os,
			t.valor_estoque_perdido,
			t.qtd_perdas,
			d.total_despesas
		FROM vendas_summary v
		CROSS JOIN os_summary o
		CROSS JOIN trocas_summary t
		CROSS JOIN despesas_summary d;
	`

	var (
		totalVendas         float64
		qtdVendas           int
		totalServicos       float64
		custoServicos       float64
		qtdOS               int
		valorEstoquePerdido float64
		qtdPerdas           int
		totalDespesas       float64
	)

	err := r.pool.QueryRow(ctx, query, inicio, fim).Scan(
		&totalVendas,
		&qtdVendas,
		&totalServicos,
		&custoServicos,
		&qtdOS,
		&valorEstoquePerdido,
		&qtdPerdas,
		&totalDespesas,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados financeiros otimizados: %w", err)
	}

	// Busca detalhamento por forma de pagamento
	formasPagamento, err := r.obterFormasPagamento(ctx, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar formas de pagamento: %w", err)
	}

	// Cálculos Derivados
	faturamentoBruto := totalVendas + totalServicos

	var ticketMedioVenda float64
	if qtdVendas > 0 {
		ticketMedioVenda = totalVendas / float64(qtdVendas)
	}

	var ticketMedioServico float64
	if qtdOS > 0 {
		ticketMedioServico = totalServicos / float64(qtdOS)
	}

	// Deduz o custo das OS no cálculo do lucro
	lucroLiquido := faturamentoBruto - custoServicos - totalDespesas - valorEstoquePerdido

	var margemLucro float64
	if faturamentoBruto > 0 {
		margemLucro = (lucroLiquido / faturamentoBruto) * 100
	}

	return &models.ResumoFinanceiro{
		PeriodoInicio:       inicio,
		PeriodoFim:          fim,
		TotalVendas:         totalVendas,
		QuantidadeVendas:    qtdVendas,
		TicketMedioVenda:    ticketMedioVenda,
		TotalServicos:       totalServicos,
		CustoServicos:       custoServicos, // Lembre-se de adicionar o campo CustoServicos no model ResumoFinanceiro se desejar exibi-lo
		QuantidadeOS:        qtdOS,
		TicketMedioServico:  ticketMedioServico,
		FaturamentoBruto:    faturamentoBruto,
		FormasPagamento:     formasPagamento,
		TotalDespesas:       totalDespesas,
		ValorEstoquePerdido: valorEstoquePerdido,
		QuantidadePerdas:    qtdPerdas,
		LucroLiquido:        lucroLiquido,
		MargemLucro:         margemLucro,
	}, nil
}

// Auxiliar para listar totais agrupados por forma de pagamento
func (r *FinanceiroRepository) obterFormasPagamento(ctx context.Context, inicio, fim time.Time) ([]models.DetalhePagamento, error) {
	query := `
		SELECT 
			COALESCE(tipo_pagamento, 'NÃO INFORMADO') AS tipo_pagamento, 
			COALESCE(SUM(total), 0) AS total
		FROM vendas
		WHERE status != 'CANCELADA'
		  AND criado_em BETWEEN $1 AND $2
		GROUP BY tipo_pagamento;
	`

	rows, err := r.pool.Query(ctx, query, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao consultar formas de pagamento: %w", err)
	}
	defer rows.Close()

	detalhes := make([]models.DetalhePagamento, 0)

	for rows.Next() {
		var d models.DetalhePagamento
		if err := rows.Scan(&d.FormaPagamento, &d.Total); err != nil {
			return nil, fmt.Errorf("erro ao ler linha de forma de pagamento: %w", err)
		}
		detalhes = append(detalhes, d)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("erro durante a iteração de formas de pagamento: %w", err)
	}

	return detalhes, nil
}

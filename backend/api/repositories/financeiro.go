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

// ObterResumoPorPeriodo agrega vendas, OS, trocas, despesas e perdas conforme o schema atual.
func (r *FinanceiroRepository) ObterResumoPorPeriodo(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	query := `
		WITH 
		vendas_summary AS (
			SELECT 
				COALESCE(SUM(total), 0) AS total_vendas,
				COALESCE(SUM(valor_troca_cedida), 0) AS valor_trocas_cedidas,
				COALESCE(SUM(desconto), 0) AS descontos_totais,
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
			  AND COALESCE(data_conclusao, criado_em) BETWEEN $1 AND $2
		),
		trocas_summary AS (
			SELECT 
				COALESCE(SUM(valor_avaliacao), 0) AS valor_trocas_recebidas,
				COUNT(*) AS qtd_trocas
			FROM trocas
			WHERE status IN ('RECEBIDA', 'AVALIADA', 'EM_ESTOQUE', 'VENDIDA')
			  AND criado_em BETWEEN $1 AND $2
		),
		despesas_summary AS (
			SELECT 
				COALESCE(SUM(valor_pago), 0) AS total_despesas
			FROM despesas
			WHERE status = 'PAGO'
			  AND data_pagamento BETWEEN $1 AND $2
		),
		perdas_summary AS (
			SELECT 
				COALESCE(SUM(custo_unitario * quantidade), 0) AS valor_estoque_perdido,
				COUNT(*) AS qtd_perdas
			FROM movimentacoes_estoque
			WHERE tipo IN ('PERDA', 'AVARIA')
			  AND criado_em BETWEEN $1 AND $2
		)
		SELECT 
			v.total_vendas,
			v.valor_trocas_cedidas,
			v.descontos_totais,
			v.qtd_vendas,
			o.total_servicos,
			o.custo_servicos,
			o.qtd_os,
			t.valor_trocas_recebidas,
			t.qtd_trocas,
			d.total_despesas,
			p.valor_estoque_perdido,
			p.qtd_perdas
		FROM vendas_summary v
		CROSS JOIN os_summary o
		CROSS JOIN trocas_summary t
		CROSS JOIN despesas_summary d
		CROSS JOIN perdas_summary p;
	`

	var (
		totalVendas          float64
		valorTrocasCedidas   float64
		descontosTotais      float64
		qtdVendas            int
		totalServicos        float64
		custoServicos        float64
		qtdOS                int
		valorTrocasRecebidas float64
		qtdTrocas            int
		totalDespesas        float64
		valorEstoquePerdido  float64
		qtdPerdas            int
	)

	err := r.pool.QueryRow(ctx, query, inicio, fim).Scan(
		&totalVendas,
		&valorTrocasCedidas,
		&descontosTotais,
		&qtdVendas,
		&totalServicos,
		&custoServicos,
		&qtdOS,
		&valorTrocasRecebidas,
		&qtdTrocas,
		&totalDespesas,
		&valorEstoquePerdido,
		&qtdPerdas,
	)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar dados financeiros otimizados: %w", err)
	}

	formasPagamento, err := r.obterFormasPagamento(ctx, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar formas de pagamento: %w", err)
	}

	receitaLiquidaVendas := totalVendas - valorTrocasCedidas
	lucroServicos := totalServicos - custoServicos
	faturamentoBruto := totalVendas + totalServicos

	var ticketMedioVenda float64
	if qtdVendas > 0 {
		ticketMedioVenda = totalVendas / float64(qtdVendas)
	}

	var ticketMedioServico float64
	if qtdOS > 0 {
		ticketMedioServico = totalServicos / float64(qtdOS)
	}

	// Lucro bruto: Faturamento - Custo Serviços - Trocas Cedidas
	lucroBruto := faturamentoBruto - custoServicos - valorTrocasCedidas
	// Lucro líquido: Lucro Bruto - Despesas + Trocas Recebidas - Perdas
	lucroLiquido := lucroBruto - totalDespesas + valorTrocasRecebidas - valorEstoquePerdido

	var margemLucro float64
	if faturamentoBruto > 0 {
		margemLucro = (lucroLiquido / faturamentoBruto) * 100
	}

	return &models.ResumoFinanceiro{
		PeriodoInicio:        inicio,
		PeriodoFim:           fim,
		TotalVendas:          totalVendas,
		ValorTrocasCedidas:   valorTrocasCedidas,
		ReceitaLiquidaVendas: receitaLiquidaVendas,
		QuantidadeVendas:     qtdVendas,
		TicketMedioVenda:     ticketMedioVenda,
		DescontosTotais:      descontosTotais,
		TotalServicos:        totalServicos,
		CustoServicos:        custoServicos,
		LucroServicos:        lucroServicos,
		QuantidadeOS:         qtdOS,
		TicketMedioServico:   ticketMedioServico,
		ValorTrocasRecebidas: valorTrocasRecebidas,
		QuantidadeTrocas:     qtdTrocas,
		FaturamentoBruto:     faturamentoBruto,
		FormasPagamento:      formasPagamento,
		TotalDespesas:        totalDespesas,
		ValorEstoquePerdido:  valorEstoquePerdido,
		QuantidadePerdas:     qtdPerdas,
		LucroBruto:           lucroBruto,
		LucroLiquido:         lucroLiquido,
		MargemLucro:          margemLucro,
	}, nil
}

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

// ObterVendasPorPeriodo lista vendas no intervalo (usado por métodos auxiliares do service).
func (r *FinanceiroRepository) ObterVendasPorPeriodo(ctx context.Context, inicio, fim time.Time) ([]models.Venda, error) {
	query := `
		SELECT 
			id, empresa_id, produto_id, usuario_id, cliente_id,
			quantidade, valor_unitario, valor_venda, valor_troca_cedida,
			desconto, total, tipo_pagamento, status, criado_em
		FROM vendas
		WHERE status != 'CANCELADA'
		  AND criado_em BETWEEN $1 AND $2
		ORDER BY criado_em DESC;
	`
	rows, err := r.pool.Query(ctx, query, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar vendas por período: %w", err)
	}
	defer rows.Close()

	vendas := make([]models.Venda, 0)
	for rows.Next() {
		var v models.Venda
		if err := rows.Scan(
			&v.ID, &v.EmpresaID, &v.ProdutoID, &v.UsuarioID, &v.ClienteID,
			&v.Quantidade, &v.ValorUnitario, &v.ValorVenda, &v.ValorTrocaCedida,
			&v.Desconto, &v.Total, &v.TipoPagamento, &v.Status, &v.CriadoEm,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler venda: %w", err)
		}
		vendas = append(vendas, v)
	}
	return vendas, rows.Err()
}

// ObterOSFinalizadasPorPeriodo lista OS finalizadas no intervalo.
func (r *FinanceiroRepository) ObterOSFinalizadasPorPeriodo(ctx context.Context, inicio, fim time.Time) ([]models.OrdemServico, error) {
	query := `
		SELECT 
			id, numero_ordem, empresa_id, usuario_id, cliente_id,
			nome_servico, defeito, servicos, valor, custo,
			assinatura, status, data_conclusao, criado_em, atualizado_em
		FROM ordens_servico
		WHERE status = 'FINALIZADA'
		  AND COALESCE(data_conclusao, criado_em) BETWEEN $1 AND $2
		ORDER BY COALESCE(data_conclusao, criado_em) DESC;
	`
	rows, err := r.pool.Query(ctx, query, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar OS finalizadas: %w", err)
	}
	defer rows.Close()

	lista := make([]models.OrdemServico, 0)
	for rows.Next() {
		var o models.OrdemServico
		if err := rows.Scan(
			&o.ID, &o.NumeroOrdem, &o.EmpresaID, &o.UsuarioID, &o.ClienteID,
			&o.NomeServico, &o.Defeito, &o.Servicos, &o.Valor, &o.Custo,
			&o.Assinatura, &o.Status, &o.DataConclusao, &o.CriadoEm, &o.AtualizadoEm,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler OS: %w", err)
		}
		lista = append(lista, o)
	}
	return lista, rows.Err()
}

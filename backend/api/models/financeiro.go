package models

import (
	"fmt"
	"time"
)

// =====================================================================
// DETALHES DE PAGAMENTO
// =====================================================================

type DetalhePagamento struct {
	FormaPagamento string  `json:"forma_pagamento"` // PIX, CARTAO_CREDITO, DINHEIRO, CHEQUE, etc
	Total          float64 `json:"total"`
}

// =====================================================================
// RESUMO FINANCEIRO COMPLETO
// =====================================================================

type ResumoFinanceiro struct {
	// Período
	PeriodoInicio time.Time `json:"periodo_inicio"`
	PeriodoFim    time.Time `json:"periodo_fim"`

	// ========== RECEITAS ==========

	// Vendas de Produtos
	TotalVendas          float64 `json:"total_vendas"`           // Valor bruto das vendas
	ValorTrocasCedidas   float64 `json:"valor_trocas_cedidas"`   // Créditos dados em trocas
	ReceitaLiquidaVendas float64 `json:"receita_liquida_vendas"` // Total - Trocas Cedidas
	QuantidadeVendas     int     `json:"quantidade_vendas"`
	TicketMedioVenda     float64 `json:"ticket_medio_venda"`
	DescontosTotais      float64 `json:"descontos_totais"`

	// Ordens de Serviço
	TotalServicos      float64 `json:"total_servicos"` // Valor de OSs finalizadas
	CustoServicos      float64 `json:"custo_servicos"` // Custo real das OSs
	LucroServicos      float64 `json:"lucro_servicos"` // Total - Custo
	QuantidadeOS       int     `json:"quantidade_os"`
	TicketMedioServico float64 `json:"ticket_medio_servico"`

	// ========== ATIVO EM ESTOQUE ==========

	// Trocas Recebidas (entrada de ativo)
	ValorTrocasRecebidas float64 `json:"valor_trocas_recebidas"` // O que entra no estoque via trocas
	QuantidadeTrocas     int     `json:"quantidade_trocas"`

	// ========== FATURAMENTO ==========

	FaturamentoBruto float64            `json:"faturamento_bruto"` // Total Vendas + Total Serviços
	FormasPagamento  []DetalhePagamento `json:"formas_pagamento"`  // Breakdown por forma

	// ========== DESPESAS ==========

	TotalDespesas float64 `json:"total_despesas"` // Despesas pagas no período

	// ========== PERDAS ==========

	ValorEstoquePerdido float64 `json:"valor_estoque_perdido"` // Perdas/avarias registradas
	QuantidadePerdas    int     `json:"quantidade_perdas"`

	// ========== RESULTADO FINAL ==========

	LucroBruto   float64 `json:"lucro_bruto"`   // Faturamento - Custo Serviços - Trocas Cedidas
	LucroLiquido float64 `json:"lucro_liquido"` // Lucro Bruto - Despesas + Trocas Recebidas
	MargemLucro  float64 `json:"margem_lucro"`  // (Lucro Líquido / Faturamento Bruto) * 100
}

// =====================================================================
// AUXILIARES
// =====================================================================

// ToString retorna uma representação em string do resumo financeiro
func (r *ResumoFinanceiro) ToString() string {
	return "ResumoFinanceiro{" +
		"PeriodoInicio: " + r.PeriodoInicio.String() +
		", PeriodoFim: " + r.PeriodoFim.String() +
		", TotalVendas: " + formatCurrency(r.TotalVendas) +
		", LucroLiquido: " + formatCurrency(r.LucroLiquido) +
		", MargemLucro: " + formatPercent(r.MargemLucro) +
		"}"
}

func formatCurrency(value float64) string {
	return fmt.Sprintf("R$ %.2f", value)
}

func formatPercent(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

// =====================================================================
// ORDENS DE SERVIÇO
// =====================================================================

// =====================================================================
// VENDAS
// =====================================================================

// =====================================================================
// TROCAS
// =====================================================================

// =====================================================================
// DESPESAS
// =====================================================================

// =====================================================================
// PRODUTOS
// =====================================================================

// =====================================================================
// CLIENTES
// =====================================================================

// =====================================================================
// MOVIMENTAÇÃO DE ESTOQUE
// =====================================================================

package models

import "time"

type DetalhePagamento struct {
	FormaPagamento string  `json:"forma_pagamento"` // PIX, CARTAO_CREDITO, DINHEIRO
	Total          float64 `json:"total"`
}

type ResumoFinanceiro struct {
	PeriodoInicio time.Time `json:"periodo_inicio"`
	PeriodoFim    time.Time `json:"periodo_fim"`

	// Vendas de Produtos
	TotalVendas      float64 `json:"total_vendas"`
	QuantidadeVendas int     `json:"quantidade_vendas"`
	TicketMedioVenda float64 `json:"ticket_medio_venda"`

	// Ordens de Serviço
	TotalServicos      float64 `json:"total_servicos"`
	QuantidadeOS       int     `json:"quantidade_os"`
	TicketMedioServico float64 `json:"ticket_medio_servico"`
	CustoServicos      float64 `json:"custo_servicos"`

	// Faturamento e Formas de Pagamento
	FaturamentoBruto float64            `json:"faturamento_bruto"`
	FormasPagamento  []DetalhePagamento `json:"formas_pagamento"`

	// Custos e Deduções
	TotalDespesas       float64 `json:"total_despesas"`
	ValorEstoquePerdido float64 `json:"valor_estoque_perdido"`
	QuantidadePerdas    int     `json:"quantidade_perdas"`

	// Indicadores
	LucroLiquido float64 `json:"lucro_liquido"`
	MargemLucro  float64 `json:"margem_lucro"` // Em porcentagem (%)
}

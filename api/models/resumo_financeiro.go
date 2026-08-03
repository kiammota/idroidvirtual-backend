package models

import "time"

type ResumoFinanceiro struct {
	PeriodoInicio time.Time `json:"periodo_inicio"`
	PeriodoFim    time.Time `json:"periodo_fim"`

	// Vendas de Produtos
	TotalVendas      float64 `json:"total_vendas"`      // R$ em vendas concluídas
	QuantidadeVendas int     `json:"quantidade_vendas"` // Quantidade de vendas realizadas

	// Ordens de Serviço
	TotalServicos float64 `json:"total_servicos"` // R$ em OS concluídas
	QuantidadeOS  int     `json:"quantidade_os"`  // Quantidade de OS finalizadas

	// Faturamento Bruto
	FaturamentoBruto float64 `json:"faturamento_bruto"` // Total Vendas + Total Serviços

	// Custos e Perdas
	TotalDespesas       float64 `json:"total_despesas"`        // Aluguel, luz, salários, etc.
	ValorEstoquePerdido float64 `json:"valor_estoque_perdido"` // R$ perdido em avarias/furtos
	QuantidadePerdas    int     `json:"quantidade_perdas"`     // Itens perdidos no período

	// Resultado Financeiro Real
	LucroLiquido float64 `json:"lucro_liquido"` // Faturamento Bruto - Despesas - Perdas
}

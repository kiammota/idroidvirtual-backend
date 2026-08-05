package models

type Fechamento struct {
	TotalVendas int     `json:"total_vendas"`
	Faturamento float64 `json:"faturamento"`
	Periodo     string  `json:"periodo"`
	Vendas      []Venda `json:"vendas"`
}

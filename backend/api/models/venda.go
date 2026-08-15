package models

import "time"

type Venda struct {
	ID               int64     `json:"id"`
	EmpresaID        int64     `json:"empresa_id"`
	ProdutoID        int64     `json:"produto_id"`
	UsuarioID        int64     `json:"usuario_id"`
	ClienteID        *int64    `json:"cliente_id"`
	Quantidade       int       `json:"quantidade"`
	ValorUnitario    float64   `json:"valor_unitario"`
	ValorVenda       float64   `json:"valor_venda"`
	ValorTrocaCedida float64   `json:"valor_troca_cedida"` // Novo campo
	Desconto         float64   `json:"desconto"`
	Total            float64   `json:"total"`
	TipoPagamento    string    `json:"tipo_pagamento"`
	Status           string    `json:"status"` // CONCLUIDA, CANCELADA
	CriadoEm         time.Time `json:"criado_em"`
}

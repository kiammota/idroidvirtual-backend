package models

import "time"

// Venda representa uma venda
type Venda struct {
	ID            int64     `json:"id"`
	ProdutoID     int64     `json:"produto_id"`
	UsuarioID     int64     `json:"usuario_id"`
	ClienteID     int64     `json:"cliente_id"`
	Quantidade    int       `json:"quantidade"`
	Total         float64   `json:"total"`
	TipoPagamento string    `json:"tipo_pagamento"`
	Status        string    `json:"status"`
	CriadoEm      time.Time `json:"criado_em"`
}

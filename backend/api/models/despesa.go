package models

import "time"

type Despesa struct {
	ID             int64      `json:"id"`
	UsuarioID      int64      `json:"usuario_id"` // Quem registou a despesa
	Descricao      string     `json:"descricao"`  // Ex: "Aluguel da Loja", "Conta de Luz", "Motoboy"
	Valor          float64    `json:"valor"`
	Categoria      string     `json:"categoria"`       // "FIXA", "VARIAVEL", "PESSOAL", "FORNECEDOR"
	Status         string     `json:"status"`          // "PENDENTE", "PAGO"
	DataVencimento time.Time  `json:"data_vencimento"` // Importante para os lembretes
	DataPagamento  *time.Time `json:"data_pagamento,omitempty"`
	CriadoEm       time.Time  `json:"criado_em"`
}

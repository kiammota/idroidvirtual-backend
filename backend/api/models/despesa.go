package models

import "time"

type Despesa struct {
	ID              int64      `json:"id"`
	EmpresaID       int64      `json:"empresa_id"`
	UsuarioID       *int64     `json:"usuario_id"`
	Descricao       string     `json:"descricao"`
	Categoria       string     `json:"categoria"`
	Valor           float64    `json:"valor"`
	ValorPago       float64    `json:"valor_pago"`
	FormaPagamento  string     `json:"forma_pagamento"`
	Status          string     `json:"status"` // PENDENTE, PAGO, ATRASADO, CANCELADO
	DataVencimento  time.Time  `json:"data_vencimento"`
	DataPagamento   *time.Time `json:"data_pagamento"`
	ComprovanteURL  string     `json:"comprovante_url"`
	NumeroDocumento string     `json:"numero_documento"`
	CriadoEm        time.Time  `json:"criado_em"`
}

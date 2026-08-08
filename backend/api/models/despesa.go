package models

import "time"

type Despesa struct {
	ID              int64      `json:"id"`
	EmpresaID       int64      `json:"empresa_id"`
	Descricao       string     `json:"descricao"`
	Categoria       string     `json:"categoria"`
	Valor           float64    `json:"valor"`
	ValorPago       float64    `json:"valor_pago"`
	FormaPagamento  *string    `json:"forma_pagamento,omitempty"`
	Status          string     `json:"status"`
	DataVencimento  time.Time  `json:"data_vencimento"`
	DataPagamento   *time.Time `json:"data_pagamento,omitempty"`
	ComprovanteURL  *string    `json:"comprovante_url,omitempty"`
	NumeroDocumento *string    `json:"numero_documento,omitempty"`
	CriadoEm        time.Time  `json:"criado_em"`
}

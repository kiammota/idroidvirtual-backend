package models

import "time"

type OrdemServico struct {
	ID           int64     `json:"id"`
	NumeroOrdem  int64     `json:"numero_ordem"`
	ProdutoID    int64     `json:"produto_id"`
	ClienteID    int64     `json:"cliente_id"`
	Defeito      *string   `json:"defeito,omitempty"`
	Servicos     *string   `json:"servicos,omitempty"`
	Assinatura   *string   `json:"assinatura,omitempty"`
	Status       string    `json:"status"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

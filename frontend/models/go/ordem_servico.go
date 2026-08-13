package models

import "time"

// OrdemServico representa uma ordem de serviço
type OrdemServico struct {
	ID          int64     `json:"id"`
	NumeroOrdem int       `json:"numero_ordem"`
	ProdutoID   int64     `json:"produto_id"`
	ClienteID   int64     `json:"cliente_id"`
	Defeito     *string   `json:"defeito"`
	Servicos    *string   `json:"servicos"`
	Assinatura  *string   `json:"assinatura"` // base64 do canvas
	Status      string    `json:"status"`
	Valor       float64   `json:"valor"`
	Custo       float64   `json:"custo"`
	CriadoEm    time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

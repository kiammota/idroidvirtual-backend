package models

import "time"

type Troca struct {
	ID             int64     `json:"id"`
	VendaID        int64     `json:"venda_id"`
	ModeloAparelho string    `json:"modelo_aparelho"`
	IMEI           string    `json:"imei"`
	ValorAvaliacao float64   `json:"valor_avaliacao"`
	Condicao       string    `json:"condicao,omitempty"`
	Observacoes    *string   `json:"observacoes,omitempty"`
	CriadoEm       time.Time `json:"criado_em"`
}

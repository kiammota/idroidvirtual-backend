package models

import "time"

// Produto representa um produto/aparelho
type Produto struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Preco        float64   `json:"preco"`
	Estoque      int       `json:"estoque"`
	Categoria    string    `json:"categoria"`
	IMEI         *string   `json:"imei,omitempty"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

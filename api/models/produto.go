package models

import "time"

type Produto struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Preco        float64   `json:"preco"`
	Estoque      int       `json:"estoque"`
	Tipo 		 string    `json:"tipo"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

package models

import "time"

type Produto struct {
	ID           int64     `json:"id"`
	EmpresaID    int64     `json:"empresa_id"`
	Nome         string    `json:"nome"`
	Preco        float64   `json:"preco"`
	Estoque      int       `json:"estoque"`
	Categoria    string    `json:"categoria"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

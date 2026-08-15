package models

import "time"

type Produto struct {
	ID           int64     `json:"id"`
	EmpresaID    int64     `json:"empresa_id"`
	Nome         string    `json:"nome"`
	Descricao    string    `json:"descricao"`
	Preco        float64   `json:"preco"`
	Categoria    string    `json:"categoria"`
	Estoque      int       `json:"estoque"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

package models

import "time"

// Cliente representa um cliente
type Cliente struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Telefone     *string   `json:"telefone"`
	Endereco     *string   `json:"endereco"`
	Email        *string   `json:"email"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

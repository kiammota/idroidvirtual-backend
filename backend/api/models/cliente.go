package models

import (
	"time"
)

type Cliente struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Telefone     *string   `json:"telefone,omitempty"`
	Endereco     *string   `json:"endereco,omitempty"`
	Email        *string   `json:"email,omitempty"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

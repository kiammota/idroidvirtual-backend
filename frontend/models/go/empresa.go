package models

import "time"

// Empresa representa uma empresa
type Empresa struct {
	ID       int64     `json:"id"`
	Nome     string    `json:"nome"`
	CNPJ     string    `json:"cnpj"`
	CriadoEm time.Time `json:"criado_em"`
}

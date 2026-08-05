package models

import "time"

// Usuario representa um usuário do sistema
type Usuario struct {
	ID        int64     `json:"id"`
	EmpresaID int64     `json:"empresa_id"`
	Email     string    `json:"email"`
	Nome      string    `json:"nome"`
	Senha     string    `json:"senha,omitempty"`
	CriadoEm  time.Time `json:"criado_em"`
}

package models

import (
	"time"
)

type Cliente struct {
	ID           int64     `json:"id"`
	EmpresaID    int64     `json:"empresa_id"`
	Nome         string    `json:"nome"`
	CPFCNPJ      string    `json:"cpf_cnpj"`
	Telefone     string    `json:"telefone"`
	Email        string    `json:"email"`
	Endereco     string    `json:"endereco"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

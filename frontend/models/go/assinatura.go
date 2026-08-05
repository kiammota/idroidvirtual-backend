package models

import "time"

// Assinatura representa assinatura de plano
type Assinatura struct {
	ID                int64      `json:"id"`
	EmpresaID         int64      `json:"empresa_id"`
	Valor             float64    `json:"valor"`
	Status            string     `json:"status"` // ATIVA, CANCELADA, PENDENTE, ATRASADA
	InicioEm          time.Time  `json:"inicio_em"`
	ProximaCobrancaEm *time.Time `json:"proxima_cobranca_em"`
	CanceladaEm       *time.Time `json:"cancelada_em"`
	CriadoEm          time.Time  `json:"criado_em"`
}

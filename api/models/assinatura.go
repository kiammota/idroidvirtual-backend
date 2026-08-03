package models

import "time"

type Assinatura struct {
	ID                int64      `json:"id"`
	EmpresaID         int64      `json:"empresa_id"`
	Valor             float64    `json:"valor"`
	Status            string     `json:"status"`
	InicioEm          time.Time  `json:"inicio_em"`
	ProximaCobrancaEm time.Time  `json:"proxima_cobranca_em"`
	CanceladaEm       *time.Time `json:"cancelada_em,omitempty"`
	CriadoEm          time.Time  `json:"criado_em"`
}

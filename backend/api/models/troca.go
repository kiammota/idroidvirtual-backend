package models

import "time"

type Troca struct {
	ID             int64     `json:"id"`
	EmpresaID      int64     `json:"empresa_id"`
	UsuarioID      int64     `json:"usuario_id"`
	ModeloAparelho string    `json:"modelo_aparelho"`
	IMEI           string    `json:"imei"`
	ValorAvaliacao float64   `json:"valor_avaliacao"`
	Condicao       string    `json:"condicao"` // USADO, SEMINOVO, EXCELENTE, DEFEITO
	Observacoes    string    `json:"observacoes"`
	Status         string    `json:"status"` // RECEBIDA, AVALIADA, EM_ESTOQUE, VENDIDA, DESCARTADA
	CriadoEm       time.Time `json:"criado_em"`
	AtualizadoEm   time.Time `json:"atualizado_em"`
}

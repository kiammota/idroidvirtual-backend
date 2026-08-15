package models

import "time"

type OrdemServico struct {
	ID              int64      `json:"id"`
	NumeroOrdem     int64      `json:"numero_ordem"`
	EmpresaID       int64      `json:"empresa_id"`
	UsuarioID       int64      `json:"usuario_id"`
	ClienteID       int64      `json:"cliente_id"`
	NomeServico     string     `json:"nome_servico"`
	Defeito         string     `json:"defeito"`
	Servicos        string     `json:"servicos"`
	Valor           float64    `json:"valor"`
	Custo           float64    `json:"custo"`
	Assinatura      string     `json:"assinatura"`
	Status          string     `json:"status"` // ABERTA, EM_ANDAMENTO, FINALIZADA, CANCELADA
	DataConclusao   *time.Time `json:"data_conclusao"`
	ClienteTelefone string     `json:"cliente_telefone,omitempty"`
	CriadoEm        time.Time  `json:"criado_em"`
	AtualizadoEm    time.Time  `json:"atualizado_em"`
}

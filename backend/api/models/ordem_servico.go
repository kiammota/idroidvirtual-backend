package models

import "time"

type OrdemServico struct {
	ID              int64     `json:"id"`
	EmpresaID       int64     `json:"empresa_id"`
	NumeroOrdem     int64     `json:"numero_ordem"`
	NomeServico     string    `json:"nome_servico"`
	UsuarioID       int64     `json:"usuario_id"`
	Valor           float64   `json:"valor"` // Geralmente float64 para representar moeda/decimais
	Custo           float64   `json:"custo"`
	ClienteID       int64     `json:"cliente_id"`
	ClienteTelefone string    `json:"cliente_telefone"`  // Util para o payload do formulário
	Cliente         *Cliente  `json:"cliente,omitempty"` // Ponteiro para trazer os dados completos do cliente
	Defeito         *string   `json:"defeito,omitempty"`
	Servicos        *string   `json:"servicos,omitempty"`
	Assinatura      *string   `json:"assinatura,omitempty"`
	Status          string    `json:"status"`
	CriadoEm        time.Time `json:"criado_em"`
	AtualizadoEm    time.Time `json:"atualizado_em"`
}

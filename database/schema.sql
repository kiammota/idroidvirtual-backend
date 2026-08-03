package models

import "time"

type Usuario struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Nome      string    `json:"nome"`
	Senha     string    `json:"senha"`
	CriadoEm  time.Time `json:"criado_em"`
}

type Produto struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Preco        float64   `json:"preco"`
	Estoque      int       `json:"estoque"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

type Cliente struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Telefone     string    `json:"telefone"`
	Endereco     string    `json:"endereco"`
	Email        string    `json:"email"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

type Venda struct {
	ID            int64     `json:"id"`
	ProdutoID     int64     `json:"produto_id"`
	Quantidade    int       `json:"quantidade"`
	Total         float64   `json:"total"`
	TipoPagamento string    `json:"tipo_pagamento"`
	CriadoEm      time.Time `json:"criado_em"`
}

type OrdemServico struct {
	ID           int64     `json:"id"`
	NumeroOrdem  int64     `json:"numero_ordem"`
	ClienteID    int64     `json:"cliente_id"`
	Defeito      string    `json:"defeito"`
	Servicos     string    `json:"servicos"`
	Assinatura   string    `json:"assinatura"`
	Status       string    `json:"status"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

type Auditoria struct {
	ID             int64                  `json:"id"`
	UsuarioID      int64                  `json:"usuario_id"`
	Acao           string                 `json:"acao"`
	Tabela         string                 `json:"tabela"`
	RegistroID     int64                  `json:"registro_id"`
	ValoresAntigos map[string]interface{} `json:"valores_antigos"`
	ValoresNovos   map[string]interface{} `json:"valores_novos"`
	Horario        time.Time              `json:"horario"`
}
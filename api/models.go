package main

import "time"

// Usuario do sistema
type Usuario struct {
	ID       int64     `json:"id"`
	Email    string    `json:"email"`
	Nome     string    `json:"nome"`
	Senha    string    `json:"-"`
	CriadoEm time.Time `json:"criado_em"`
}

// Produto em estoque
type Produto struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Preco        float64   `json:"preco"`
	Estoque      int       `json:"estoque"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Cliente
type Cliente struct {
	ID           int64     `json:"id"`
	Nome         string    `json:"nome"`
	Telefone     string    `json:"telefone"`
	Endereco     string    `json:"endereco"`
	Email        string    `json:"email"`
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Venda realizada
type Venda struct {
	ID            int64     `json:"id"`
	ProdutoID     int64     `json:"produto_id"`
	Quantidade    int       `json:"quantidade"`
	Total         float64   `json:"total"`
	TipoPagamento string    `json:"tipo_pagamento"` // cartao, maquina, pix
	CriadoEm      time.Time `json:"criado_em"`
}

// Ordem de Serviço
type OrdemServico struct {
	ID           int64     `json:"id"`
	NumeroOrdem  int64     `json:"numero_ordem"`
	ClienteID    int64     `json:"cliente_id"`
	Defeito      string    `json:"defeito"`
	Servicos     string    `json:"servicos"`
	Assinatura   string    `json:"assinatura"`
	Status       string    `json:"status"` // pendente, em_andamento, pronto, entregue
	CriadoEm     time.Time `json:"criado_em"`
	AtualizadoEm time.Time `json:"atualizado_em"`
}

// Dashboard - resumo do dia
type Painel struct {
	VendasHoje     float64 `json:"vendas_hoje"`
	TotalProdutos  int     `json:"total_produtos"`
	OrdensEmAberto int     `json:"ordens_em_aberto"`
	EstoquesBaixos int     `json:"estoques_baixos"`
}

// Requests

type RegistroRequest struct {
	Email string `json:"email" binding:"required"`
	Nome  string `json:"nome" binding:"required"`
	Senha string `json:"senha" binding:"required"`
}

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Senha string `json:"senha" binding:"required"`
}

type CriarProdutoRequest struct {
	Nome    string  `json:"nome" binding:"required"`
	Preco   float64 `json:"preco" binding:"required"`
	Estoque int     `json:"estoque" binding:"required"`
}

type AtualizarProdutoRequest struct {
	Nome    string  `json:"nome"`
	Preco   float64 `json:"preco"`
	Estoque int     `json:"estoque"`
}

type CriarClienteRequest struct {
	Nome     string `json:"nome" binding:"required"`
	Telefone string `json:"telefone" binding:"required"`
	Endereco string `json:"endereco" binding:"required"`
	Email    string `json:"email"`
}

type AtualizarClienteRequest struct {
	Nome     string `json:"nome"`
	Telefone string `json:"telefone"`
	Endereco string `json:"endereco"`
	Email    string `json:"email"`
}

type CriarVendaRequest struct {
	ProdutoID     int64  `json:"produto_id" binding:"required"`
	Quantidade    int    `json:"quantidade" binding:"required"`
	TipoPagamento string `json:"tipo_pagamento" binding:"required"`
}

type CriarOrdemServicoRequest struct {
	ClienteID int64  `json:"cliente_id" binding:"required"`
	Defeito   string `json:"defeito" binding:"required"`
	Servicos  string `json:"servicos" binding:"required"`
}

type AtualizarStatusOrdemRequest struct {
	Status string `json:"status" binding:"required"`
}

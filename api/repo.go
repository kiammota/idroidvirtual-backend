package main

import (
	"context"
	"log"
)

// ============ USUARIOS ============

func BuscarUsuarioPorEmail(email string) (*Usuario, error) {
	var usuario Usuario
	var senha_hash string
	err := bd.QueryRow(context.Background(),
		"SELECT id, email, nome, senha, criado_em FROM usuarios WHERE email = $1",
		email).Scan(&usuario.ID, &usuario.Email, &usuario.Nome, &senha_hash, &usuario.CriadoEm)

	if err != nil {
		return nil, err
	}

	usuario.Senha = senha_hash
	return &usuario, nil
}

func CriarUsuario(email, nome, senha_hash string) (int64, error) {
	var usuario_id int64
	err := bd.QueryRow(context.Background(),
		"INSERT INTO usuarios (email, nome, senha) VALUES ($1, $2, $3) RETURNING id",
		email, nome, senha_hash).Scan(&usuario_id)

	if err != nil {
		return 0, err
	}

	return usuario_id, nil
}

func BuscarUsuarioPorID(id int64) (*Usuario, error) {
	var usuario Usuario
	err := bd.QueryRow(context.Background(),
		"SELECT id, email, nome, criado_em FROM usuarios WHERE id = $1",
		id).Scan(&usuario.ID, &usuario.Email, &usuario.Nome, &usuario.CriadoEm)

	if err != nil {
		return nil, err
	}

	return &usuario, nil
}

// ============ PRODUTOS ============

func CriarProduto(nome string, preco float64, estoque int) (int64, error) {
	var id int64
	err := bd.QueryRow(context.Background(),
		"INSERT INTO produtos (nome, preco, estoque) VALUES ($1, $2, $3) RETURNING id",
		nome, preco, estoque).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func ListarProdutos() ([]Produto, error) {
	linhas, err := bd.Query(context.Background(),
		"SELECT id, nome, preco, estoque, criado_em, atualizado_em FROM produtos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var produtos []Produto
	for linhas.Next() {
		var p Produto
		err := linhas.Scan(&p.ID, &p.Nome, &p.Preco, &p.Estoque, &p.CriadoEm, &p.AtualizadoEm)
		if err != nil {
			log.Println(err)
			continue
		}
		produtos = append(produtos, p)
	}

	return produtos, nil
}

func BuscarProdutoPorID(id int64) (*Produto, error) {
	var p Produto
	err := bd.QueryRow(context.Background(),
		"SELECT id, nome, preco, estoque, criado_em, atualizado_em FROM produtos WHERE id = $1",
		id).Scan(&p.ID, &p.Nome, &p.Preco, &p.Estoque, &p.CriadoEm, &p.AtualizadoEm)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func AtualizarProduto(id int64, nome string, preco float64, estoque int) error {
	_, err := bd.Exec(context.Background(),
		"UPDATE produtos SET nome = $1, preco = $2, estoque = $3, atualizado_em = NOW() WHERE id = $4",
		nome, preco, estoque, id)

	return err
}

func DeletarProduto(id int64) error {
	_, err := bd.Exec(context.Background(),
		"DELETE FROM produtos WHERE id = $1", id)

	return err
}

// ============ CLIENTES ============

func CriarCliente(nome, telefone, endereco, email string) (int64, error) {
	var id int64
	err := bd.QueryRow(context.Background(),
		"INSERT INTO clientes (nome, telefone, endereco, email) VALUES ($1, $2, $3, $4) RETURNING id",
		nome, telefone, endereco, email).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func ListarClientes() ([]Cliente, error) {
	linhas, err := bd.Query(context.Background(),
		"SELECT id, nome, telefone, endereco, email, criado_em, atualizado_em FROM clientes ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var clientes []Cliente
	for linhas.Next() {
		var c Cliente
		err := linhas.Scan(&c.ID, &c.Nome, &c.Telefone, &c.Endereco, &c.Email, &c.CriadoEm, &c.AtualizadoEm)
		if err != nil {
			log.Println(err)
			continue
		}
		clientes = append(clientes, c)
	}

	return clientes, nil
}

func BuscarClientePorID(id int64) (*Cliente, error) {
	var c Cliente
	err := bd.QueryRow(context.Background(),
		"SELECT id, nome, telefone, endereco, email, criado_em, atualizado_em FROM clientes WHERE id = $1",
		id).Scan(&c.ID, &c.Nome, &c.Telefone, &c.Endereco, &c.Email, &c.CriadoEm, &c.AtualizadoEm)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func AtualizarCliente(id int64, nome, telefone, endereco, email string) error {
	_, err := bd.Exec(context.Background(),
		"UPDATE clientes SET nome = $1, telefone = $2, endereco = $3, email = $4, atualizado_em = NOW() WHERE id = $5",
		nome, telefone, endereco, email, id)

	return err
}

func DeletarCliente(id int64) error {
	_, err := bd.Exec(context.Background(),
		"DELETE FROM clientes WHERE id = $1", id)

	return err
}

// ============ VENDAS ============

func CriarVenda(produto_id int64, quantidade int, total float64, tipo_pagamento string) (int64, error) {
	var venda_id int64
	err := bd.QueryRow(context.Background(),
		"INSERT INTO vendas (produto_id, quantidade, total, tipo_pagamento) VALUES ($1, $2, $3, $4) RETURNING id",
		produto_id, quantidade, total, tipo_pagamento).Scan(&venda_id)

	if err != nil {
		return 0, err
	}

	return venda_id, nil
}

func ListarVendas() ([]Venda, error) {
	linhas, err := bd.Query(context.Background(),
		"SELECT id, produto_id, quantidade, total, tipo_pagamento, criado_em FROM vendas ORDER BY criado_em DESC")
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var vendas []Venda
	for linhas.Next() {
		var v Venda
		err := linhas.Scan(&v.ID, &v.ProdutoID, &v.Quantidade, &v.Total, &v.TipoPagamento, &v.CriadoEm)
		if err != nil {
			log.Println(err)
			continue
		}
		vendas = append(vendas, v)
	}

	return vendas, nil
}

func DecrementarEstoque(produto_id int64, quantidade int) error {
	_, err := bd.Exec(context.Background(),
		"UPDATE produtos SET estoque = estoque - $1 WHERE id = $2",
		quantidade, produto_id)

	return err
}

// ============ ORDENS DE SERVIÇO ============

func ObterProximoNumeroOrdem() (int64, error) {
	var proximo_numero int64
	err := bd.QueryRow(context.Background(),
		"SELECT COALESCE(MAX(numero_ordem), 0) + 1 FROM ordens_servico").Scan(&proximo_numero)

	if err != nil {
		return 0, err
	}

	return proximo_numero, nil
}

func CriarOrdemServico(numero_ordem int64, cliente_id int64, defeito, servicos string) (int64, error) {
	var id int64
	err := bd.QueryRow(context.Background(),
		"INSERT INTO ordens_servico (numero_ordem, cliente_id, defeito, servicos, status) VALUES ($1, $2, $3, $4, 'pendente') RETURNING id",
		numero_ordem, cliente_id, defeito, servicos).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func ListarOrdensServico() ([]OrdemServico, error) {
	linhas, err := bd.Query(context.Background(),
		"SELECT id, numero_ordem, cliente_id, defeito, servicos, status, criado_em, atualizado_em FROM ordens_servico ORDER BY criado_em DESC")
	if err != nil {
		return nil, err
	}
	defer linhas.Close()

	var ordens []OrdemServico
	for linhas.Next() {
		var o OrdemServico
		err := linhas.Scan(&o.ID, &o.NumeroOrdem, &o.ClienteID, &o.Defeito, &o.Servicos, &o.Status, &o.CriadoEm, &o.AtualizadoEm)
		if err != nil {
			log.Println(err)
			continue
		}
		ordens = append(ordens, o)
	}

	return ordens, nil
}

func AtualizarStatusOrdem(id int64, status string) error {
	_, err := bd.Exec(context.Background(),
		"UPDATE ordens_servico SET status = $1, atualizado_em = NOW() WHERE id = $2",
		status, id)

	return err
}

// ============ DASHBOARD ============

func ObterResumoVendas() (float64, error) {
	var vendas_hoje float64
	err := bd.QueryRow(context.Background(),
		"SELECT COALESCE(SUM(total), 0) FROM vendas WHERE DATE(criado_em) = CURRENT_DATE").Scan(&vendas_hoje)

	if err != nil {
		return 0, err
	}

	return vendas_hoje, nil
}

func ContarProdutos() (int, error) {
	var total int
	err := bd.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM produtos").Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func ContarOrdensEmAberto() (int, error) {
	var total int
	err := bd.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM ordens_servico WHERE status IN ('pendente', 'em_andamento')").Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func ContarEstoquesBaixos() (int, error) {
	var total int
	err := bd.QueryRow(context.Background(),
		"SELECT COUNT(*) FROM produtos WHERE estoque < 5").Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

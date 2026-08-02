package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ============ AUTENTICAÇÃO ============

func AutenticarHandler(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	usuario, err := BuscarUsuarioPorEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "email ou senha inválidos"})
		return
	}

	if !VerificarSenha(usuario.Senha, req.Senha) {
		c.JSON(http.StatusUnauthorized, gin.H{"erro": "email ou senha inválidos"})
		return
	}

	token, err := GerarToken(usuario.ID, usuario.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao gerar token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "usuario": usuario})
}

func RegistrarHandler(c *gin.Context) {
	var req RegistroRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	senha_hash, err := GerarHashSenha(req.Senha)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao processar senha"})
		return
	}

	usuario_id, err := CriarUsuario(req.Email, req.Nome, senha_hash)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "usuário já existe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"mensagem": "usuário criado com sucesso", "usuario_id": usuario_id})
}

func ObterMeuUsuarioHandler(c *gin.Context) {
	usuario_id := ObterUsuarioIDDoContexto(c)

	usuario, err := BuscarUsuarioPorID(usuario_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "usuário não encontrado"})
		return
	}

	c.JSON(http.StatusOK, usuario)
}

// ============ PRODUTOS ============

func CriarProdutoHandler(c *gin.Context) {
	var req CriarProdutoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	id, err := CriarProduto(req.Nome, req.Preco, req.Estoque)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao criar produto"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "mensagem": "produto criado com sucesso"})
}

func ListarProdutosHandler(c *gin.Context) {
	produtos, err := ListarProdutos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao listar produtos"})
		return
	}

	if produtos == nil {
		produtos = []Produto{}
	}

	c.JSON(http.StatusOK, produtos)
}

func AtualizarProdutoHandler(c *gin.Context) {
	var req AtualizarProdutoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	produto_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	err = AtualizarProduto(produto_id, req.Nome, req.Preco, req.Estoque)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao atualizar produto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "produto atualizado com sucesso"})
}

func DeletarProdutoHandler(c *gin.Context) {
	produto_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	err = DeletarProduto(produto_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao deletar produto"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "produto deletado com sucesso"})
}

// ============ CLIENTES ============

func CriarClienteHandler(c *gin.Context) {
	var req CriarClienteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	id, err := CriarCliente(req.Nome, req.Telefone, req.Endereco, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao criar cliente"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "mensagem": "cliente criado com sucesso"})
}

func ListarClientesHandler(c *gin.Context) {
	clientes, err := ListarClientes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao listar clientes"})
		return
	}

	if clientes == nil {
		clientes = []Cliente{}
	}

	c.JSON(http.StatusOK, clientes)
}

func AtualizarClienteHandler(c *gin.Context) {
	var req AtualizarClienteRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	cliente_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	err = AtualizarCliente(cliente_id, req.Nome, req.Telefone, req.Endereco, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao atualizar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "cliente atualizado com sucesso"})
}

func DeletarClienteHandler(c *gin.Context) {
	cliente_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	err = DeletarCliente(cliente_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao deletar cliente"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "cliente deletado com sucesso"})
}

// ============ VENDAS ============

func CriarVendaHandler(c *gin.Context) {
	var req CriarVendaRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	produto, err := BuscarProdutoPorID(req.ProdutoID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "produto não encontrado"})
		return
	}

	if produto.Estoque < req.Quantidade {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "estoque insuficiente"})
		return
	}

	total := produto.Preco * float64(req.Quantidade)

	venda_id, err := CriarVenda(req.ProdutoID, req.Quantidade, total, req.TipoPagamento)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao criar venda"})
		return
	}

	err = DecrementarEstoque(req.ProdutoID, req.Quantidade)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao atualizar estoque"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": venda_id, "total": total, "mensagem": "venda registrada com sucesso"})
}

func ListarVendasHandler(c *gin.Context) {
	vendas, err := ListarVendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao listar vendas"})
		return
	}

	if vendas == nil {
		vendas = []Venda{}
	}

	c.JSON(http.StatusOK, vendas)
}

// ============ ORDENS DE SERVIÇO ============

func CriarOrdemServicoHandler(c *gin.Context) {
	var req CriarOrdemServicoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Verificar se cliente existe
	_, err := BuscarClientePorID(req.ClienteID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erro": "cliente não encontrado"})
		return
	}

	numero_ordem, err := ObterProximoNumeroOrdem()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao gerar número da ordem"})
		return
	}

	id, err := CriarOrdemServico(numero_ordem, req.ClienteID, req.Defeito, req.Servicos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao criar ordem de serviço"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "numero_ordem": numero_ordem, "mensagem": "ordem de serviço criada com sucesso"})
}

func ListarOrdensServicoHandler(c *gin.Context) {
	ordens, err := ListarOrdensServico()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao listar ordens de serviço"})
		return
	}

	if ordens == nil {
		ordens = []OrdemServico{}
	}

	c.JSON(http.StatusOK, ordens)
}

func AtualizarStatusOrdemHandler(c *gin.Context) {
	var req AtualizarStatusOrdemRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	ordem_id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "id inválido"})
		return
	}

	err = AtualizarStatusOrdem(ordem_id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao atualizar status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mensagem": "status atualizado com sucesso"})
}

// ============ DASHBOARD ============

func ObterPainelHandler(c *gin.Context) {
	vendas_hoje, err := ObterResumoVendas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao obter vendas"})
		return
	}

	total_produtos, err := ContarProdutos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao contar produtos"})
		return
	}

	ordens_em_aberto, err := ContarOrdensEmAberto()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao contar ordens"})
		return
	}

	estoques_baixos, err := ContarEstoquesBaixos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": "falha ao contar estoques"})
		return
	}

	painel := Painel{
		VendasHoje:     vendas_hoje,
		TotalProdutos:  total_produtos,
		OrdensEmAberto: ordens_em_aberto,
		EstoquesBaixos: estoques_baixos,
	}

	c.JSON(http.StatusOK, painel)
}

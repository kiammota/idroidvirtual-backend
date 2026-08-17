package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/models"
	"backend/repositories"
)

type VendaService struct {
	pool        *pgxpool.Pool
	auditoria   AuditoriaService
	vendaRepo   repositories.VendaRepository
	usuarioRepo repositories.UsuarioRepository
	clienteRepo repositories.ClienteRepository
	produtoRepo repositories.ProdutoRepository
}

func NewVendaService(
	pool *pgxpool.Pool,
	vendaRepo repositories.VendaRepository,
	produtoRepo repositories.ProdutoRepository,
	usuarioRep repositories.UsuarioRepository,
	clienteRep repositories.ClienteRepository,
	audit AuditoriaService,
) *VendaService {
	return &VendaService{
		pool:        pool,
		vendaRepo:   vendaRepo,
		produtoRepo: produtoRepo,
		usuarioRepo: usuarioRep,
		clienteRepo: clienteRep,
		auditoria:   audit,
	}
}

func (s *VendaService) ListarPorPeriodo(ctx context.Context, empresaID int64, inicio time.Time, fim time.Time) ([]models.Venda, error) {
	if inicio.After(fim) {
		return nil, errors.New("a data inicial não pode ser posterior à data final")
	}

	vendas, err := s.vendaRepo.BuscarPorPeriodo(ctx, empresaID, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar vendas por período: %w", err)
	}

	return vendas, nil
}

func (s *VendaService) ObterVendasDoDia(ctx context.Context, empresaID int64) ([]models.Venda, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())

	return s.vendaRepo.BuscarPorPeriodo(ctx, empresaID, inicioDia, fimDia)
}

func (s *VendaService) ObterVendasDoMes(ctx context.Context, empresaId int64) ([]models.Venda, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fimMes := inicioMes.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.vendaRepo.BuscarPorPeriodo(ctx, empresaId, inicioMes, fimMes)
}

// ExecutarVenda realiza a venda, baixa o estoque e registra na auditoria.
// clienteID = 0 significa Consumidor Final (opcional).
func (s *VendaService) ExecutarVenda(
	ctx context.Context,
	usuarioID int64,
	produtoID int64,
	quantidade int,
	tipoPagamento string,
	clienteID int64, // 0 = Consumidor Final
) (*models.Venda, error) {

	if quantidade <= 0 {
		return nil, errors.New("a quantidade deve ser maior que zero")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// 1. Produto
	produto, err := s.produtoRepo.BuscarPorIDTx(ctx, tx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	if produto == nil {
		return nil, errors.New("produto não encontrado")
	}
	if produto.Estoque < quantidade {
		if produto.Estoque == 0 {
			return nil, errors.New("O estoque para este produto está zerado!")
		}
		return nil, fmt.Errorf(
			"O Estoque está insuficiente! (disponível: %d, solicitado: %d)",
			produto.Estoque, quantidade,
		)
	}

	// 2. Usuário
	usuario, err := s.usuarioRepo.BuscarPorIDTx(ctx, tx, usuarioID)
	if err != nil || usuario == nil {
		return nil, errors.New("usuário associado à venda não foi encontrado")
	}

	// 3. Cliente (opcional)
	if clienteID > 0 {
		cliente, err := s.clienteRepo.BuscarPorIDTx(ctx, tx, clienteID)
		if err != nil {
			return nil, fmt.Errorf("erro ao buscar cliente: %w", err)
		}
		if cliente == nil {
			return nil, errors.New("cliente informado não foi encontrado")
		}
	}

	// 4. Baixa de estoque
	if err := s.produtoRepo.DiminuirEstoqueTx(ctx, tx, produtoID, quantidade); err != nil {
		return nil, fmt.Errorf("falha ao dar baixa no estoque: %w", err)
	}

	valorTotal := produto.Preco * float64(quantidade)

	venda := &models.Venda{
		ProdutoID:     produtoID,
		UsuarioID:     usuarioID,
		ClienteID:     clienteID, // 0 = Consumidor Final
		Quantidade:    quantidade,
		Total:         valorTotal,
		TipoPagamento: tipoPagamento,
		Status:        "CONCLUIDA",
	}

	if err := s.vendaRepo.CriarTx(ctx, tx, venda); err != nil {
		return nil, fmt.Errorf("falha ao registrar venda: %w", err)
	}

	err = s.auditoria.CadastrarVenda(ctx, tx, usuarioID, venda.ID, valorTotal, tipoPagamento)
	if err != nil {
		return nil, fmt.Errorf("falha ao gravar auditoria da venda: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("erro ao efetivar transação: %w", err)
	}

	return venda, nil
}

// CancelarVenda estorna o estoque do produto e cancela a venda atomicamente
func (s *VendaService) CancelarVenda(
	ctx context.Context,
	usuarioID int64,
	vendaID int64,
) error {

	venda, err := s.vendaRepo.BuscarPorID(ctx, vendaID)
	if err != nil {
		return fmt.Errorf("erro ao buscar venda: %w", err)
	}
	if venda == nil {
		return errors.New("venda não encontrada")
	}

	if venda.Status == "CANCELADA" {
		return errors.New("esta venda já está cancelada")
	}

	produto, err := s.produtoRepo.BuscarPorID(ctx, int(venda.ProdutoID))
	if err != nil || produto == nil {
		return errors.New("produto associado à venda não foi encontrado")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.produtoRepo.AumentarEstoqueTx(ctx, tx, venda.ProdutoID, venda.Quantidade)
	if err != nil {
		return fmt.Errorf("falha ao devolver produto ao estoque: %w", err)
	}

	err = s.vendaRepo.CancelarTx(ctx, tx, vendaID)
	if err != nil {
		return fmt.Errorf("falha ao cancelar registro da venda: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao efetivar cancelamento: %w", err)
	}

	return nil
}

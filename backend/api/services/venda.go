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
	vendaRepo   repositories.VendaRepository
	produtoRepo repositories.ProdutoRepository
	auditRepo   repositories.AuditoriaRepository
}

func NewVendaService(
	pool *pgxpool.Pool,
	vendaRepo repositories.VendaRepository,
	produtoRepo repositories.ProdutoRepository,
	auditRepo repositories.AuditoriaRepository,
) *VendaService {
	return &VendaService{
		pool:        pool,
		vendaRepo:   vendaRepo,
		produtoRepo: produtoRepo,
		auditRepo:   auditRepo,
	}
}

func (s *VendaService) ObterVendasDoDia(ctx context.Context) ([]models.Venda, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())

	return s.vendaRepo.BuscarPorPeriodo(ctx, inicioDia, fimDia)
}

func (s *VendaService) ObterVendasDoMes(ctx context.Context) ([]models.Venda, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fimMes := inicioMes.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.vendaRepo.BuscarPorPeriodo(ctx, inicioMes, fimMes)
}

// ExecutarVenda realiza a venda, baixa o estoque e registra na auditoria
func (s *VendaService) ExecutarVenda(
	ctx context.Context,
	usuarioID int64,
	produtoID int64,
	quantidade int,
	tipoPagamento string,
) (*models.Venda, error) {

	if quantidade <= 0 {
		return nil, errors.New("a quantidade deve ser maior que zero")
	}

	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	produto, err := s.produtoRepo.BuscarPorIDTx(ctx, tx, produtoID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produto: %w", err)
	}
	if produto == nil {
		return nil, errors.New("produto não encontrado")
	}

	if produto.Estoque < quantidade {
		return nil, fmt.Errorf(
			"estoque insuficiente (disponível: %d, solicitado: %d)",
			produto.Estoque,
			quantidade,
		)
	}

	if err := s.produtoRepo.DiminuirEstoqueTx(ctx, tx, produtoID, quantidade); err != nil {
		return nil, fmt.Errorf("falha ao dar baixa no estoque: %w", err)
	}

	valorTotal := produto.Preco * float64(quantidade)

	venda := &models.Venda{
		ProdutoID:     produtoID,
		UsuarioID:     usuarioID,
		Quantidade:    quantidade,
		Total:         valorTotal,
		TipoPagamento: tipoPagamento,
		Status:        "CONCLUIDA",
	}

	if err := s.vendaRepo.CriarTx(ctx, tx, venda); err != nil {
		return nil, fmt.Errorf("falha ao registrar venda: %w", err)
	}

	// Se você tiver um model de auditoria, registre aqui.
	// Exemplo:
	//
	// auditoria := &models.Auditoria{
	// 	UsuarioID: usuarioID,
	// 	Acao:      "VENDA_CRIADA",
	// 	Detalhes:  fmt.Sprintf("Venda do produto %d, quantidade %d", produtoID, quantidade),
	// }
	//
	// if err := s.auditRepo.CriarTx(ctx, tx, auditoria); err != nil {
	// 	return nil, fmt.Errorf("falha ao gravar auditoria: %w", err)
	// }

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
	motivo string,
) error {

	// 1. Busca a Venda original
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

	// 2. Busca o Produto atual para salvar o estoque antes do estorno
	produto, err := s.produtoRepo.BuscarPorID(ctx, venda.ProdutoID)
	if err != nil || produto == nil {
		return errors.New("produto associado à venda não foi encontrado")
	}

	// 3. Inicia a Transação
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// 4. Estorna a quantidade de volta para o Estoque
	err = s.produtoRepo.AumentarEstoqueTx(ctx, tx, venda.ProdutoID, venda.Quantidade)
	if err != nil {
		return fmt.Errorf("falha ao devover produto ao estoque: %w", err)
	}

	// 5. Atualiza o status da venda para CANCELADA
	err = s.vendaRepo.CancelarTx(ctx, tx, vendaID)
	if err != nil {
		return fmt.Errorf("falha ao cancelar registro da venda: %w", err)
	}

	if err != nil {
		return fmt.Errorf("falha ao gravar auditoria do cancelamento: %w", err)
	}

	// 7. Efetiva todas as operações
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao efetivar cancelamento: %w", err)
	}

	return nil
}

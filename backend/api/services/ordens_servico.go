package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"backend/models"
	"backend/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdensServicoService struct {
	Pool              *pgxpool.Pool
	OrdensServicoRepo *repositories.OrdemServicoRepository
	ClienteRepo       *repositories.ClienteRepository
	AuditService      *AuditoriaService
}

func NewOrdensServicoService(
	pool *pgxpool.Pool,
	ordemRepo *repositories.OrdemServicoRepository,
	clienteRepo *repositories.ClienteRepository,
	auditService *AuditoriaService,
) *OrdensServicoService {
	return &OrdensServicoService{
		Pool:              pool,
		OrdensServicoRepo: ordemRepo,
		ClienteRepo:       clienteRepo,
		AuditService:      auditService,
	}
}

func (s *OrdensServicoService) ListarPorPeriodo(ctx context.Context, empresaID int64, inicio, fim time.Time) ([]models.OrdemServico, error) {
	if inicio.After(fim) {
		return nil, errors.New("a data inicial não pode ser maior que a data final")
	}

	return s.OrdensServicoRepo.ListarPorPeriodo(ctx, empresaID, inicio, fim)
}

// POST /service-orders
func (s *OrdensServicoService) Criar(ctx context.Context, ordem *models.OrdemServico) (err error) {
	if ordem == nil {
		return errors.New("ordem de serviço inválida")
	}

	if ordem.Status == "" {
		ordem.Status = "ABERTA"
	}

	ordem.NomeServico = strings.ToUpper(strings.TrimSpace(ordem.NomeServico))
	if ordem.NomeServico == "" {
		return errors.New("nome do serviço é obrigatório")
	}

	// 1. Extrai o telefone de forma segura (sem dar panic se ordem.Cliente for nil)
	var telefoneBruto string
	if ordem.ClienteTelefone != "" {
		telefoneBruto = ordem.ClienteTelefone
	} else if ordem.Cliente != nil && ordem.Cliente.Telefone != nil {
		telefoneBruto = *ordem.Cliente.Telefone
	}

	telefoneLimpo := apenasNumeros(telefoneBruto)
	if telefoneLimpo == "" {
		return errors.New("telefone do cliente é obrigatório")
	}

	// 2. Inicia transação
	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %w", err)
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// 3. Busca cliente existente ou cria um novo
	cliente, err := s.ClienteRepo.BuscarPorTelefoneTx(ctx, tx, telefoneLimpo)
	if err != nil {
		return fmt.Errorf("erro ao buscar cliente: %w", err)
	}

	if cliente == nil {
		enderecoPadrao := "NÃO INFORMADO"
		novoCliente := &models.Cliente{
			Nome:     "NÃO INFORMADO",
			Telefone: &telefoneLimpo,
			Endereco: &enderecoPadrao,
		}

		if err = s.ClienteRepo.CriarTx(ctx, tx, novoCliente); err != nil {
			return fmt.Errorf("erro ao criar cliente: %w", err)
		}

		cliente = novoCliente
	}

	// 4. Vincula o ID do cliente encontrado/criado na ordem
	ordem.ClienteID = cliente.ID

	// 5. Salva a Ordem de Serviço
	if err = s.OrdensServicoRepo.CriarTx(ctx, tx, ordem); err != nil {
		return fmt.Errorf("erro ao salvar ordem de serviço: %w", err)
	}

	// 6. Registro de Auditoria
	if err = s.AuditService.CadastrarOrdemDeServico(ctx, tx, ordem.UsuarioID, ordem.ID, cliente.Nome); err != nil {
		return fmt.Errorf("erro ao registrar auditoria: %w", err)
	}

	// 7. Confirma transação
	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return nil
}

// Função auxiliar para remover máscara/caracteres do telefone
func apenasNumeros(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// GET /service-orders
func (s *OrdensServicoService) Listar(ctx context.Context, empresaID int64, limit, offset int) ([]models.OrdemServico, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.OrdensServicoRepo.Listar(ctx, empresaID, limit, offset)
}

func (s *OrdensServicoService) ListarPorStatus(ctx context.Context, status string, limit, offset int) ([]models.OrdemServico, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	statusLimpo := strings.ToUpper(strings.TrimSpace(status))

	if statusLimpo != "EM_ANDAMENTO" && statusLimpo != "FINALIZADA" && statusLimpo != "CANCELADA" {
		return nil, errors.New("status inválido: informe EM_ANDAMENTO, FINALIZADA ou CANCELADA")
	}

	return s.OrdensServicoRepo.ListarPorStatus(ctx, statusLimpo, limit, offset)
}

// GET /service-orders/{id}
func (s *OrdensServicoService) BuscarPorID(ctx context.Context, id int64) (*models.OrdemServico, error) {
	if id <= 0 {
		return nil, errors.New("id inválido")
	}

	ordem, err := s.OrdensServicoRepo.BuscarPorID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ordem == nil {
		return nil, errors.New("ordem de serviço não encontrada")
	}

	return ordem, nil
}

// PATCH /service-orders/{id}
func (s *OrdensServicoService) Atualizar(ctx context.Context, id int64, patch *models.OrdemServico) (err error) {
	if id <= 0 {
		return errors.New("id inválido")
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	existente, err := s.OrdensServicoRepo.BuscarPorIDTx(ctx, tx, id)
	if err != nil {
		return err
	}
	if existente == nil {
		return errors.New("ordem de serviço não encontrada")
	}

	atualizado := *existente

	if patch.NumeroOrdem != 0 {
		if patch.NumeroOrdem <= 0 {
			return errors.New("número de ordem inválido")
		}
		atualizado.NumeroOrdem = patch.NumeroOrdem
	}

	if patch.ClienteID != 0 {
		if patch.ClienteID <= 0 {
			return errors.New("cliente inválido")
		}
		cliente, err := s.ClienteRepo.BuscarPorIDTx(ctx, tx, patch.ClienteID)
		if err != nil {
			return err
		}
		if cliente == nil {
			return errors.New("cliente não encontrado")
		}
		atualizado.ClienteID = patch.ClienteID
	}

	if patch.Defeito != nil {
		atualizado.Defeito = patch.Defeito
	}

	if patch.Servicos != nil {
		atualizado.Servicos = patch.Servicos
	}

	if patch.Valor != 0 {
		atualizado.Valor = patch.Valor
	}

	if patch.Custo != 0 {
		atualizado.Custo = patch.Custo
	}

	if patch.Assinatura != nil {
		atualizado.Assinatura = patch.Assinatura
	}

	if patch.Status != "" {
		atualizado.Status = patch.Status
	}

	atualizado.ID = id

	if err = s.OrdensServicoRepo.AtualizarTx(ctx, tx, &atualizado); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// DELETE /service-orders/{id}
func (s *OrdensServicoService) Deletar(ctx context.Context, id int64) (err error) {
	if id <= 0 {
		return errors.New("id inválido")
	}

	tx, err := s.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	ordem, err := s.OrdensServicoRepo.BuscarPorIDTx(ctx, tx, id)
	if err != nil {
		return err
	}
	if ordem == nil {
		return errors.New("ordem de serviço não encontrada")
	}

	if err = s.OrdensServicoRepo.DeletarTx(ctx, tx, id); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GET /service-orders/client-id/{cliente_id}
func (s *OrdensServicoService) BuscarPorClienteID(ctx context.Context, clienteID int64, limit, offset int) ([]models.OrdemServico, error) {
	if clienteID <= 0 {
		return nil, errors.New("cliente inválido")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.OrdensServicoRepo.BuscarPorClienteID(ctx, clienteID, limit, offset)
}

// GET /service-orders/client-name/{client_name}
func (s *OrdensServicoService) BuscarPorNomeCliente(ctx context.Context, nome string, limit, offset int) ([]models.OrdemServico, error) {
	if nome == "" {
		return nil, errors.New("nome do cliente obrigatório")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.OrdensServicoRepo.BuscarPorNomeCliente(ctx, nome, limit, offset)
}

// GET /service-orders/client-email/{client_email}
func (s *OrdensServicoService) BuscarPorEmailCliente(ctx context.Context, email string, limit, offset int) ([]models.OrdemServico, error) {
	if email == "" {
		return nil, errors.New("email obrigatório")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.OrdensServicoRepo.BuscarPorEmailCliente(ctx, email, limit, offset)
}

// GET /service-orders/client-number/{client_number}
func (s *OrdensServicoService) BuscarPorNumeroCliente(ctx context.Context, numero string, limit, offset int) ([]models.OrdemServico, error) {
	if numero == "" {
		return nil, errors.New("número obrigatório")
	}
	if limit <= 0 {
		limit = 10
	}
	return s.OrdensServicoRepo.BuscarPorNumeroCliente(ctx, numero, limit, offset)
}

// GET /service-orders/three-months
func (s *OrdensServicoService) UltimosTresMeses(ctx context.Context, empresaID int64) ([]models.OrdemServico, error) {
	return s.OrdensServicoRepo.UltimosTresMeses(ctx, empresaID)
}

// GET /service-orders/six-months
func (s *OrdensServicoService) UltimosSeisMeses(ctx context.Context, empresaID int64) ([]models.OrdemServico, error) {
	return s.OrdensServicoRepo.UltimosSeisMeses(ctx, empresaID)
}

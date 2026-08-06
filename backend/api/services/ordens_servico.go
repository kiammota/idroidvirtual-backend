package services

import (
	"context"
	"errors"
	"strings"

	"backend/models"
	"backend/repositories"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdemServicoPatch struct {
	NumeroOrdem *int64  `json:"numero_ordem,omitempty"`
	ClienteID   *int64  `json:"cliente_id,omitempty"`
	Defeito     *string `json:"defeito,omitempty"`
	Servicos    *string `json:"servicos,omitempty"`
	Assinatura  *string `json:"assinatura,omitempty"`
	Status      *string `json:"status,omitempty"`
}

type OrdensServicoService struct {
	Pool              *pgxpool.Pool
	OrdensServicoRepo *repositories.OrdemServicoRepository
	ClienteRepo       *repositories.ClienteRepository
}

func NewOrdensServicoService(
	pool *pgxpool.Pool,
	ordemRepo *repositories.OrdemServicoRepository,
	clienteRepo *repositories.ClienteRepository,
) *OrdensServicoService {
	return &OrdensServicoService{
		Pool:              pool,
		OrdensServicoRepo: ordemRepo,
		ClienteRepo:       clienteRepo,
	}
}

// POST /service-orders
func (s *OrdensServicoService) Criar(ctx context.Context, ordem *models.OrdemServico) (err error) {
	if ordem.Status == "" {
		ordem.Status = "ABERTA"
	}

	// Força o nome do serviço a ficar em caixa alta
	ordem.NomeServico = strings.ToUpper(strings.TrimSpace(ordem.NomeServico))
	if ordem.NomeServico == "" {
		return errors.New("nome do serviço é obrigatório")
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

	var cliente *models.Cliente

	// 1. Tenta buscar o cliente existente pelo telefone
	telefoneLimpo := apenasNumeros(ordem.ClienteTelefone)
	if telefoneLimpo != "" {
		cliente, err = s.ClienteRepo.BuscarPorTelefoneTx(ctx, tx, telefoneLimpo)
		if err != nil {
			return err
		}
	}

	// 2. Se não encontrou o cliente pelo telefone, cria um novo
	if cliente == nil {

		// Prepara o ponteiro para o Telefone (se a string não for vazia)
		var telPtr *string
		if telefoneLimpo != "" {
			telPtr = &telefoneLimpo
		}

		// Prepara o ponteiro para o Endereço padrão
		enderecoPadrao := "NÃO INFORMADO"

		novoCliente := &models.Cliente{
			Nome:     "NÃO INFORMADO",
			Telefone: telPtr,
			Endereco: &enderecoPadrao,
		}

		// Cria o cliente no banco dentro da transação (retornando o ID gerado)
		if err = s.ClienteRepo.CriarTx(ctx, tx, novoCliente); err != nil {
			return err
		}

		cliente = novoCliente
	}

	// 3. Pega o ID do cliente (seja ele existente ou recém-criado) e insere na ordem
	ordem.ClienteID = cliente.ID

	// 4. Busca a quantidade atual de ordens em andamento e gera o número da OS
	emAndamento, err := s.OrdensServicoRepo.ContarEmAndamentoTx(ctx, tx)
	if err != nil {
		return err
	}
	ordem.NumeroOrdem = int64(emAndamento + 1)

	// 5. Salva a Ordem de Serviço com o ClienteID preenchido
	if err = s.OrdensServicoRepo.CriarTx(ctx, tx, ordem); err != nil {
		return err
	}

	return tx.Commit(ctx)
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
func (s *OrdensServicoService) Listar(ctx context.Context, limit, offset int) ([]models.OrdemServico, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}
	return s.OrdensServicoRepo.Listar(ctx, limit, offset)
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
func (s *OrdensServicoService) Atualizar(ctx context.Context, id int64, patch *OrdemServicoPatch) (err error) {
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

	if patch.NumeroOrdem != nil {
		if *patch.NumeroOrdem <= 0 {
			return errors.New("número de ordem inválido")
		}
		atualizado.NumeroOrdem = *patch.NumeroOrdem
	}

	if patch.ClienteID != nil {
		if *patch.ClienteID <= 0 {
			return errors.New("cliente inválido")
		}
		cliente, err := s.ClienteRepo.BuscarPorIDTx(ctx, tx, *patch.ClienteID)
		if err != nil {
			return err
		}
		if cliente == nil {
			return errors.New("cliente não encontrado")
		}
		atualizado.ClienteID = *patch.ClienteID
	}

	if patch.Defeito != nil {
		atualizado.Defeito = patch.Defeito
	}

	if patch.Servicos != nil {
		atualizado.Servicos = patch.Servicos
	}

	if patch.Assinatura != nil {
		atualizado.Assinatura = patch.Assinatura
	}

	if patch.Status != nil {
		atualizado.Status = *patch.Status
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
func (s *OrdensServicoService) UltimosTresMeses(ctx context.Context) ([]models.OrdemServico, error) {
	return s.OrdensServicoRepo.UltimosTresMeses(ctx)
}

// GET /service-orders/six-months
func (s *OrdensServicoService) UltimosSeisMeses(ctx context.Context) ([]models.OrdemServico, error) {
	return s.OrdensServicoRepo.UltimosSeisMeses(ctx)
}

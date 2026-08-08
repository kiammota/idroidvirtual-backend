package services

import (
	"context"
	"fmt"
	"time"

	"backend/models"
	"backend/repositories"

	"github.com/jackc/pgx/v5"
)

type PeriodoTipo string

const (
	PeriodoDia    PeriodoTipo = "dia"
	PeriodoSemana PeriodoTipo = "semana"
	PeriodoMes    PeriodoTipo = "mes"
)

type AuditoriaService struct {
	repo        *repositories.AuditoriaRepository
	empresaRepo *repositories.EmpresaRepository
}

func NewAuditoriaService(repo *repositories.AuditoriaRepository, empresaRepo *repositories.EmpresaRepository) *AuditoriaService {
	return &AuditoriaService{
		repo:        repo,
		empresaRepo: empresaRepo, // Corrigido: agora armazena a referência
	}
}

// log registra ações genéricas de auditoria obtendo o empresaID dinamicamente
func (s *AuditoriaService) log(
	ctx context.Context,
	usuarioID int64,
	acao string,
) error {
	empresa, err := s.empresaRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return fmt.Errorf("falha ao identificar empresa do usuário para auditoria: %w", err)
	}

	audit := models.Auditoria{
		EmpresaID: empresa.ID,
		UsuarioID: usuarioID,
		Acao:      acao,
	}
	return s.repo.Criar(ctx, audit)
}

// logTx registra ações de auditoria dentro de uma transação ativa
func (s *AuditoriaService) logTx(
	ctx context.Context,
	tx pgx.Tx,
	usuarioID int64,
	acao string,
) error {
	empresa, err := s.empresaRepo.BuscarPorUsuarioID(ctx, usuarioID)
	if err != nil {
		return fmt.Errorf("falha ao identificar empresa do usuário para auditoria: %w", err)
	}

	audit := models.Auditoria{
		EmpresaID: empresa.ID,
		UsuarioID: usuarioID,
		Acao:      acao,
	}
	return s.repo.CriarTx(ctx, tx, audit)
}

// Métodos auxiliares específicos para registrar ações comuns:

func (s *AuditoriaService) CadastrarVenda(
	ctx context.Context,
	tx pgx.Tx,
	usuarioID int64,
	vendaID int64,
	valor float64,
	tipoPagamento string,
) error {
	acao := fmt.Sprintf("Cadastrou a venda #%d no valor de R$ %.2f via %s", vendaID, valor, tipoPagamento)
	if tx != nil {
		return s.logTx(ctx, tx, usuarioID, acao)
	}
	return s.log(ctx, usuarioID, acao)
}

func (s *AuditoriaService) CadastrarOrdemDeServico(ctx context.Context, tx pgx.Tx, usuarioID int64, osID int64, cliente string) error {
	acao := fmt.Sprintf("Cadastrou a ordem de serviço #%d para o cliente %s", osID, cliente)
	if tx != nil {
		return s.logTx(ctx, tx, usuarioID, acao)
	}
	return s.log(ctx, usuarioID, acao)
}

func (s *AuditoriaService) CadastrarUsuario(ctx context.Context, tx pgx.Tx, usuarioID int64, novoUsuarioNome string) error {
	acao := fmt.Sprintf("Cadastrou o novo usuário: %s", novoUsuarioNome)
	if tx != nil {
		return s.logTx(ctx, tx, usuarioID, acao)
	}
	return s.log(ctx, usuarioID, acao)
}

// ListarPorPeriodo busca os registros considerando o fuso de São Paulo
func (s *AuditoriaService) ListarPorPeriodo(
	ctx context.Context,
	periodo PeriodoTipo,
	pagina int,
	itensPorPagina int,
) ([]models.Auditoria, error) {

	if pagina < 1 {
		pagina = 1
	}
	if itensPorPagina <= 0 {
		itensPorPagina = 20
	}

	offset := (pagina - 1) * itensPorPagina

	// Define a timezone de São Paulo para cálculo dos períodos
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.FixedZone("UTC-3", -3*60*60)
	}

	agora := time.Now().In(loc)
	var inicio time.Time
	fim := agora

	switch periodo {
	case PeriodoDia:
		inicio = time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, loc)
	case PeriodoSemana:
		inicio = agora.AddDate(0, 0, -7)
	case PeriodoMes:
		inicio = time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, loc)
	default:
		inicio = time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, loc)
	}

	return s.repo.BuscarPorPeriodo(ctx, inicio, fim, itensPorPagina, offset)
}

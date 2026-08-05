package services

import (
	"context"
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
	repo *repositories.AuditoriaRepository
}

func NewAuditoriaService(repo *repositories.AuditoriaRepository) *AuditoriaService {
	return &AuditoriaService{repo: repo}
}

// Log registra ações de auditoria fora de uma transação
func (s *AuditoriaService) Log(
	ctx context.Context,
	usuarioID int64,
	acao string,
	tabela string,
	registroID int64,
	antigo map[string]interface{},
	novo map[string]interface{},
) error {
	audit := models.Auditoria{
		UsuarioID:      usuarioID,
		Acao:           acao,
		Tabela:         tabela,
		RegistroID:     registroID,
		ValoresAntigos: antigo,
		ValoresNovos:   novo,
	}
	return s.repo.Criar(ctx, audit)
}

// LogTx registra a auditoria reaproveitando uma transação ativa (pgx.Tx)
func (s *AuditoriaService) LogTx(
	ctx context.Context,
	tx pgx.Tx,
	usuarioID int64,
	acao string,
	tabela string,
	registroID int64,
	antigo map[string]interface{},
	novo map[string]interface{},
) error {
	audit := models.Auditoria{
		UsuarioID:      usuarioID,
		Acao:           acao,
		Tabela:         tabela,
		RegistroID:     registroID,
		ValoresAntigos: antigo,
		ValoresNovos:   novo,
	}
	return s.repo.CriarTx(ctx, tx, audit)
}

// ListarPorPeriodo busca os registros calculando o intervalo de tempo (dia, semana, mês)
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
	agora := time.Now()
	var inicio time.Time
	fim := agora

	switch periodo {
	case PeriodoDia:
		inicio = time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	case PeriodoSemana:
		inicio = agora.AddDate(0, 0, -7)
	case PeriodoMes:
		inicio = time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	default:
		inicio = time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	}

	return s.repo.BuscarPorPeriodo(ctx, inicio, fim, itensPorPagina, offset)
}

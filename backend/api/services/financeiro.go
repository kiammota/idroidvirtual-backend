package services

import (
	"context"
	"fmt"
	"time"

	"backend/models"
	"backend/repositories"
)

type FinanceiroService struct {
	financeiroRepo *repositories.FinanceiroRepository
}

func NewFinanceiroService(financeiroRepo *repositories.FinanceiroRepository) *FinanceiroService {
	return &FinanceiroService{
		financeiroRepo: financeiroRepo,
	}
}

// FechamentoPorPeriodo busca o resumo do dia atual (00:00:00 até 23:59:59)
func (s *FinanceiroService) FechamentoPorPeriodo(ctx context.Context, empresaId int64) (*models.ResumoFinanceiro, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())

	return s.ObterFinanceiroPorIntervalo(ctx, empresaId, inicioDia, fimDia)
}

// ObterFinanceiroDoMes calcula do primeiro dia do mês até o instante atual
func (s *FinanceiroService) ObterFinanceiroDoMes(ctx context.Context, empresaID int64) (*models.ResumoFinanceiro, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fimMes := inicioMes.AddDate(0, 1, 0).Add(-time.Nanosecond)

	return s.ObterFinanceiroPorIntervalo(ctx, empresaID, inicioMes, fimMes)
}

// ObterFinanceiroPorIntervalo valida o período de datas e delega a busca otimizada ao repositório
func (s *FinanceiroService) ObterFinanceiroPorIntervalo(ctx context.Context, empresaID int64, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	if inicio.After(fim) {
		return nil, fmt.Errorf("a data inicial (%s) não pode ser posterior à data final (%s)", inicio.Format("2006-01-02"), fim.Format("2006-01-02"))
	}

	resumo, err := s.financeiroRepo.ObterResumoPorPeriodo(ctx, empresaID, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular resumo financeiro: %w", err)
	}

	return resumo, nil
}

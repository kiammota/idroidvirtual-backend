package services

// import (
// 	"context"
// 	"fmt"
// 	"time"

// 	"backend/models"
// 	"backend/repositories"
// )

// type FinanceiroService struct {
// 	vendaRepo repositories.VendaRepository
// 	osRepo    repositories.OrdemServicoRepository
// }

// func NewFinanceiroService(
// 	vendaRepo repositories.VendaRepository,
// 	osRepo repositories.OrdemServicoRepository,
// ) *FinanceiroService {
// 	return &FinanceiroService{
// 		vendaRepo: vendaRepo,
// 		osRepo:    osRepo,
// 	}
// }

// // ObterBalancoDiario calcula o caixa do dia atual para o fechamento de turno
// func (s *FinanceiroService) ObterBalancoDiario(ctx context.Context) (*models.ResumoFinanceiro, error) {
// 	agora := time.Now()
// 	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
// 	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())

// 	return s.ObterBalancoPorPeriodo(ctx, inicioDia, fimDia)
// }

// // ObterBalancoPorPeriodo busca os dados consolidados sem alterar nada no banco
// func (s *FinanceiroService) ObterBalancoPorPeriodo(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
// 	// 1. Busca total de Vendas
// 	totalVendas, err := s.vendaRepo.ObterFaturamentoPorPeriodo(ctx, inicio, fim)
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao buscar faturamento de vendas: %w", err)
// 	}

// 	// 2. Busca total de Despesas Pagas
// 	if err != nil {
// 		return nil, fmt.Errorf("erro ao buscar despesas: %w", err)
// 	}

// 	// 3. Monta o DTO de retorno

// 	resumo := &models.ResumoFinanceiro{
// 		PeriodoInicio: inicio,
// 		PeriodoFim:    fim,
// 		TotalVendas:   totalVendas,
// 		TotalDespesas: totalDespesas,
// 		LucroLiquido:  lucro,
// 	}

// 	return resumo, nil
// }

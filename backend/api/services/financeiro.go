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

// FechamentoDia - Resumo do dia atual
func (s *FinanceiroService) FechamentoDia(ctx context.Context) (*models.ResumoFinanceiro, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.FechamentoPorIntervalo(ctx, inicioDia, fimDia)
}

// FechamentoPorPeriodo é alias usado pelo controller para o fechamento do dia.
func (s *FinanceiroService) FechamentoPorPeriodo(ctx context.Context) (*models.ResumoFinanceiro, error) {
	return s.FechamentoDia(ctx)
}

// FechamentoMes - Resumo do mês atual (do dia 1 até agora)
func (s *FinanceiroService) FechamentoMes(ctx context.Context) (*models.ResumoFinanceiro, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fim := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.FechamentoPorIntervalo(ctx, inicioMes, fim)
}

// ObterFinanceiroDoMes alias para o controller.
func (s *FinanceiroService) ObterFinanceiroDoMes(ctx context.Context) (*models.ResumoFinanceiro, error) {
	return s.FechamentoMes(ctx)
}

// FechamentoMesCompleto - Resumo do mês passado inteiro
func (s *FinanceiroService) FechamentoMesCompleto(ctx context.Context) (*models.ResumoFinanceiro, error) {
	agora := time.Now()
	mesPassado := agora.AddDate(0, -1, 0)
	inicioMes := time.Date(mesPassado.Year(), mesPassado.Month(), 1, 0, 0, 0, 0, mesPassado.Location())
	fimMes := inicioMes.AddDate(0, 1, 0).Add(-time.Nanosecond)
	return s.FechamentoPorIntervalo(ctx, inicioMes, fimMes)
}

// FechamentoPersonalizado / FechamentoPorIntervalo
func (s *FinanceiroService) FechamentoPersonalizado(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	return s.FechamentoPorIntervalo(ctx, inicio, fim)
}

func (s *FinanceiroService) FechamentoPorIntervalo(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	if ctx == nil {
		return nil, fmt.Errorf("contexto não pode ser nulo")
	}
	if inicio.After(fim) {
		return nil, fmt.Errorf(
			"a data inicial (%s) não pode ser posterior à data final (%s)",
			inicio.Format("2006-01-02"),
			fim.Format("2006-01-02"),
		)
	}
	resumo, err := s.financeiroRepo.ObterResumoPorPeriodo(ctx, inicio, fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao calcular resumo financeiro: %w", err)
	}
	return resumo, nil
}

// ObterFinanceiroPorIntervalo alias para o controller.
func (s *FinanceiroService) ObterFinanceiroPorIntervalo(ctx context.Context, inicio, fim time.Time) (*models.ResumoFinanceiro, error) {
	return s.FechamentoPorIntervalo(ctx, inicio, fim)
}

func (s *FinanceiroService) ObterVendasDoDia(ctx context.Context) ([]models.Venda, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.financeiroRepo.ObterVendasPorPeriodo(ctx, inicioDia, fimDia)
}

func (s *FinanceiroService) ObterVendasDoMes(ctx context.Context) ([]models.Venda, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fimMes := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.financeiroRepo.ObterVendasPorPeriodo(ctx, inicioMes, fimMes)
}

func (s *FinanceiroService) ObterOSFinalizadasDoDia(ctx context.Context) ([]models.OrdemServico, error) {
	agora := time.Now()
	inicioDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, agora.Location())
	fimDia := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.financeiroRepo.ObterOSFinalizadasPorPeriodo(ctx, inicioDia, fimDia)
}

func (s *FinanceiroService) ObterOSFinalizadasDoMes(ctx context.Context) ([]models.OrdemServico, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())
	fimMes := time.Date(agora.Year(), agora.Month(), agora.Day(), 23, 59, 59, 999999999, agora.Location())
	return s.financeiroRepo.ObterOSFinalizadasPorPeriodo(ctx, inicioMes, fimMes)
}

func (s *FinanceiroService) FormatarResumo(resumo *models.ResumoFinanceiro) string {
	if resumo == nil {
		return "Resumo nulo"
	}
	return fmt.Sprintf(`
╔════════════════════════════════════════════════════════════╗
║                    RESUMO FINANCEIRO                       ║
╠════════════════════════════════════════════════════════════╣
║ Período: %s até %s
║
║ RECEITAS
║ ─────────────────────────────────────────────────────────
║ Vendas (total):           R$ %.2f
║ Trocas cedidas:           R$ %.2f
║ Receita líquida (vendas): R$ %.2f
║ Serviços:                 R$ %.2f
║ Lucro em serviços:        R$ %.2f
║
║ ATIVO (Estoque)
║ ─────────────────────────────────────────────────────────
║ Valor de trocas recebidas: R$ %.2f
║
║ DESPESAS
║ ─────────────────────────────────────────────────────────
║ Total de despesas:        R$ %.2f
║ Perdas/avarias:           R$ %.2f
║
║ RESULTADO FINAL
║ ─────────────────────────────────────────────────────────
║ Faturamento bruto:        R$ %.2f
║ Lucro bruto:              R$ %.2f
║ Lucro líquido:            R$ %.2f
║ Margem de lucro:          %.2f%%
║
║ QUANTIDADE
║ ─────────────────────────────────────────────────────────
║ Vendas:                   %d
║ Ordens de serviço:        %d
║ Trocas recebidas:         %d
║ Registros de perda:       %d
╚════════════════════════════════════════════════════════════╝
	`,
		resumo.PeriodoInicio.Format("02/01/2006"),
		resumo.PeriodoFim.Format("02/01/2006"),
		resumo.TotalVendas,
		resumo.ValorTrocasCedidas,
		resumo.ReceitaLiquidaVendas,
		resumo.TotalServicos,
		resumo.LucroServicos,
		resumo.ValorTrocasRecebidas,
		resumo.TotalDespesas,
		resumo.ValorEstoquePerdido,
		resumo.FaturamentoBruto,
		resumo.LucroBruto,
		resumo.LucroLiquido,
		resumo.MargemLucro,
		resumo.QuantidadeVendas,
		resumo.QuantidadeOS,
		resumo.QuantidadeTrocas,
		resumo.QuantidadePerdas,
	)
}

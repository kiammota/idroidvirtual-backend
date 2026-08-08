package services

import (
	"context"
	"fmt"
	"time"

	"backend/models"
	"backend/repositories"
)

type DespesaService struct {
	repo *repositories.DespesaRepository
}

func NewDespesaService(repo *repositories.DespesaRepository) *DespesaService {
	return &DespesaService{repo: repo}
}

func (s *DespesaService) Criar(ctx context.Context, d *models.Despesa) error {
	if d.Descricao == "" {
		return fmt.Errorf("a descrição é obrigatória")
	}
	if d.Valor <= 0 {
		return fmt.Errorf("o valor da despesa deve ser maior que zero")
	}

	// Regra de negócio: Ajusta status padrão se não enviado
	if d.Status == "" {
		d.Status = "PENDENTE"
	}

	// Se o valor pago for igual ou maior que o valor total, marca como PAGO
	if d.ValorPago >= d.Valor {
		d.Status = "PAGO"
		if d.DataPagamento == nil {
			agora := time.Now()
			d.DataPagamento = &agora
		}
	}

	return s.repo.Criar(ctx, d)
}

func (s *DespesaService) BuscarPorID(ctx context.Context, id, empresaID int64) (*models.Despesa, error) {
	return s.repo.BuscarPorID(ctx, id, empresaID)
}

func (s *DespesaService) ListarPorPeriodo(ctx context.Context, empresaID int64, inicio, fim time.Time) ([]models.Despesa, error) {
	if inicio.After(fim) {
		return nil, fmt.Errorf("data inicial não pode ser maior que a data final")
	}
	return s.repo.ListarPorEmpresaEPeriodo(ctx, empresaID, inicio, fim)
}

func (s *DespesaService) Atualizar(ctx context.Context, d *models.Despesa) error {
	if d.ID == 0 {
		return fmt.Errorf("id da despesa inválido")
	}

	// Atualiza status baseado no valor pago
	if d.ValorPago >= d.Valor {
		d.Status = "PAGO"
	}

	return s.repo.Atualizar(ctx, d)
}

func (s *DespesaService) Deletar(ctx context.Context, id, empresaID int64) error {
	return s.repo.Deletar(ctx, id, empresaID)
}

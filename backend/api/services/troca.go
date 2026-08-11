package services

import (
	"context"
	"errors"
	"fmt"

	"backend/models"
	"backend/repositories"
)

// TrocaService define o contrato dos métodos do serviço de troca
type TrocaService interface {
	ProcessarTroca(ctx context.Context, troca *models.Troca) error
	ObterPorVendaID(ctx context.Context, vendaID int64) (*models.Troca, error)
	Listar(ctx context.Context, limit, offset int) ([]*models.Troca, error)
	ListarUltimoMes(ctx context.Context, limit, offset int) ([]*models.Troca, error)
	ListarUltimosTresMeses(ctx context.Context, limit, offset int) ([]*models.Troca, error)
	ListarUltimosSeisMeses(ctx context.Context, limit, offset int) ([]*models.Troca, error)
}

type trocaService struct {
	trocaRepo repositories.TrocaRepository
}

func NewTrocaService(
	trocaRepo repositories.TrocaRepository,
) TrocaService {
	return &trocaService{
		trocaRepo: trocaRepo,
	}
}

func (s *trocaService) ProcessarTroca(
	ctx context.Context,
	troca *models.Troca,
) error {

	if troca.VendaID <= 0 {
		return errors.New("o ID da venda associada é obrigatório")
	}

	if troca.ModeloAparelho == "" {
		return errors.New("o modelo do aparelho de troca é obrigatório")
	}

	if troca.IMEI == "" {
		return errors.New("o IMEI do aparelho de troca é obrigatório")
	}

	if troca.ValorAvaliacao <= 0 {
		return errors.New("o valor de avaliação da troca deve ser maior que zero")
	}

	if err := s.trocaRepo.Criar(ctx, troca); err != nil {
		return fmt.Errorf(
			"erro ao registrar troca: %w",
			err,
		)
	}

	return nil
}

func (s *trocaService) ObterPorVendaID(
	ctx context.Context,
	vendaID int64,
) (*models.Troca, error) {

	if vendaID <= 0 {
		return nil, errors.New("ID da venda inválido")
	}

	troca, err := s.trocaRepo.BuscarPorVendaID(
		ctx,
		vendaID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao buscar troca: %w",
			err,
		)
	}

	return troca, nil
}

func (s *trocaService) Listar(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {

	trocas, err := s.trocaRepo.Listar(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas: %w", err)
	}

	return trocas, nil
}

func (s *trocaService) ListarUltimoMes(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {

	trocas, err := s.trocaRepo.ListarUltimoMes(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas do último mês: %w", err)
	}

	return trocas, nil
}

func (s *trocaService) ListarUltimosTresMeses(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {

	trocas, err := s.trocaRepo.ListarUltimosTresMeses(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas dos últimos 3 meses: %w", err)
	}

	return trocas, nil
}

func (s *trocaService) ListarUltimosSeisMeses(
	ctx context.Context,
	limit, offset int,
) ([]*models.Troca, error) {

	trocas, err := s.trocaRepo.ListarUltimosSeisMeses(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas dos últimos 6 meses: %w", err)
	}

	return trocas, nil
}

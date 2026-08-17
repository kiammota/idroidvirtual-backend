package services

import (
	"context"
	"errors"
	"fmt"

	"backend/models"
	"backend/repositories"
)

type TrocaService struct {
	trocaRepo repositories.TrocaRepository
}

func NewTrocaService(
	trocaRepo repositories.TrocaRepository,
) TrocaService {
	return TrocaService{
		trocaRepo: trocaRepo,
	}
}

func (s *TrocaService) ProcessarTroca(
	ctx context.Context,
	troca *models.Troca,
) error {

	if troca.EmpresaID <= 0 {
		return errors.New("o ID da empresa é obrigatório")
	}

	if troca.ModeloAparelho == "" {
		return errors.New("o modelo do aparelho de troca é obrigatório")
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

func (s *TrocaService) ObterPorID(
	ctx context.Context,
	id int64,
	empresaID int64,
) (*models.Troca, error) {

	if id <= 0 {
		return nil, errors.New("ID inválido")
	}

	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}

	troca, err := s.trocaRepo.BuscarPorID(
		ctx,
		id,
		empresaID,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"erro ao buscar troca: %w",
			err,
		)
	}

	return troca, nil
}

func (s *TrocaService) Listar(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {

	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}

	trocas, err := s.trocaRepo.Listar(ctx, empresaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas: %w", err)
	}

	return trocas, nil
}

func (s *TrocaService) ListarUltimoMes(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {

	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}

	trocas, err := s.trocaRepo.ListarUltimoMes(ctx, empresaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas do último mês: %w", err)
	}

	return trocas, nil
}

func (s *TrocaService) ListarUltimosTresMeses(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {

	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}

	trocas, err := s.trocaRepo.ListarUltimosTresMeses(ctx, empresaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas dos últimos 3 meses: %w", err)
	}

	return trocas, nil
}

func (s *TrocaService) ListarUltimosSeisMeses(
	ctx context.Context,
	empresaID int64,
	limit, offset int,
) ([]*models.Troca, error) {

	if empresaID <= 0 {
		return nil, errors.New("ID da empresa inválido")
	}

	trocas, err := s.trocaRepo.ListarUltimosSeisMeses(ctx, empresaID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("erro ao listar trocas dos últimos 6 meses: %w", err)
	}

	return trocas, nil
}

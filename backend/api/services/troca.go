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

func (s *TrocaService) ObterPorVendaID(
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

package repositories

import (
	"backend/models"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdemServicoRepository struct {
	Pool *pgxpool.Pool
}

func NewOrdemServicoRepository(pool *pgxpool.Pool) OrdemServicoRepository {
	return OrdemServicoRepository{Pool: pool}
}

// Criar insere uma nova ordem de serviço e retorna o ID gerado e as datas do banco
func (r *OrdemServicoRepository) Criar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		INSERT INTO ordens_servico (
			numero_ordem, cliente_id, defeito, servicos, assinatura, status
		) VALUES (
			$1, $2, $3, $4, $5, $6
		)
		RETURNING id, criado_em, atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		os.NumeroOrdem,
		os.ClienteID,
		os.Defeito,
		os.Servicos,
		os.Assinatura,
		os.Status,
	).Scan(&os.ID, &os.CriadoEm, &os.AtualizadoEm)

	if err != nil {
		return err
	}

	return nil
}

// BuscarPorID retorna uma ordem de serviço pelo seu ID
func (r *OrdemServicoRepository) BuscarPorID(ctx context.Context, id int64) (*models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE id = $1
	`

	var os models.OrdemServico
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&os.ID,
		&os.NumeroOrdem,
		&os.ClienteID,
		&os.Defeito,
		&os.Servicos,
		&os.Assinatura,
		&os.Status,
		&os.CriadoEm,
		&os.AtualizadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Retorna nil se não encontrar registro
		}
		return nil, err
	}

	return &os, nil
}

// Listar retorna todas as ordens de serviço (com paginação simples de limite/offset)
func (r *OrdemServicoRepository) Listar(ctx context.Context, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		ORDER BY id DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordens []models.OrdemServico
	for rows.Next() {
		var os models.OrdemServico
		err := rows.Scan(
			&os.ID,
			&os.NumeroOrdem,
			&os.ClienteID,
			&os.Defeito,
			&os.Servicos,
			&os.Assinatura,
			&os.Status,
			&os.CriadoEm,
			&os.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		ordens = append(ordens, os)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ordens, nil
}

// ============================================================================
// MÉTODOS DE FILTRO TEMPORAL
// ============================================================================

// ListarPorPeriodo busca ordens de serviço criadas entre duas datas específicas
func (r *OrdemServicoRepository) ListarPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE criado_em >= $1 AND criado_em <= $2
		ORDER BY criado_em DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.Pool.Query(ctx, query, dataInicio, dataFim, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ordens []models.OrdemServico
	for rows.Next() {
		var os models.OrdemServico
		err := rows.Scan(
			&os.ID,
			&os.NumeroOrdem,
			&os.ClienteID,
			&os.Defeito,
			&os.Servicos,
			&os.Assinatura,
			&os.Status,
			&os.CriadoEm,
			&os.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		ordens = append(ordens, os)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ordens, nil
}

// ListarDesteMes busca as ordens de serviço criadas a partir do dia 1º do mês atual
func (r *OrdemServicoRepository) ListarDesteMes(ctx context.Context, limit, offset int) ([]models.OrdemServico, error) {
	agora := time.Now()
	inicioMes := time.Date(agora.Year(), agora.Month(), 1, 0, 0, 0, 0, agora.Location())

	return r.ListarPorPeriodo(ctx, inicioMes, agora, limit, offset)
}

// ListarUltimosTresMeses busca as ordens de serviço dos últimos 3 meses até o instante atual
func (r *OrdemServicoRepository) ListarUltimosTresMeses(ctx context.Context, limit, offset int) ([]models.OrdemServico, error) {
	agora := time.Now()
	tresMesesAtras := agora.AddDate(0, -3, 0)

	return r.ListarPorPeriodo(ctx, tresMesesAtras, agora, limit, offset)
}

func (r *OrdemServicoRepository) ListarDoAnoAtual(ctx context.Context, limit, offset int) ([]models.OrdemServico, error) {
	agora := time.Now()
	inicioAno := time.Date(agora.Year(), time.January, 1, 0, 0, 0, 0, agora.Location())

	return r.ListarPorPeriodo(ctx, inicioAno, agora, limit, offset)
}

// ============================================================================

// Atualizar altera os campos de uma ordem de serviço existente
func (r *OrdemServicoRepository) Atualizar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		UPDATE ordens_servico
		SET numero_ordem = $1,
		    cliente_id = $2,
		    defeito = $3,
		    servicos = $4,
		    assinatura = $5,
		    status = $6,
		    atualizado_em = NOW()
		WHERE id = $7
		RETURNING atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		os.NumeroOrdem,
		os.ClienteID,
		os.Defeito,
		os.Servicos,
		os.Assinatura,
		os.Status,
		os.ID,
	).Scan(&os.AtualizadoEm)

	if err != nil {
		return err
	}

	return nil
}

// Finalizar altera o status da ordem de serviço para "FINALIZADA"
func (r *OrdemServicoRepository) Finalizar(ctx context.Context, id int64) error {
	query := `
		UPDATE ordens_servico
		SET status = 'FINALIZADA',
		    atualizado_em = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhuma ordem de serviço encontrada para finalizar")
	}

	return nil
}

package repositories

import (
	"context"
	"errors"
	"time"

	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrdemServicoRepository struct {
	Pool *pgxpool.Pool
}

func NewOrdemServicoRepository(pool *pgxpool.Pool) *OrdemServicoRepository {
	return &OrdemServicoRepository{Pool: pool}
}

// ============================================================================
// MÉTODOS DE ESCRITA / MUTACIONAL
// ============================================================================

func (r *OrdemServicoRepository) Criar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		INSERT INTO ordens_servico (numero_ordem, cliente_id, defeito, servicos, assinatura, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, criado_em, atualizado_em
	`
	return r.Pool.QueryRow(
		ctx, query,
		os.NumeroOrdem, os.ClienteID, os.Defeito, os.Servicos, os.Assinatura, os.Status,
	).Scan(&os.ID, &os.CriadoEm, &os.AtualizadoEm)
}

func (r *OrdemServicoRepository) CriarTx(ctx context.Context, tx pgx.Tx, os *models.OrdemServico) error {
	query := `
		INSERT INTO ordens_servico (numero_ordem, cliente_id, defeito, servicos, assinatura, status)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, criado_em, atualizado_em
	`
	return tx.QueryRow(
		ctx, query,
		os.NumeroOrdem, os.ClienteID, os.Defeito, os.Servicos, os.Assinatura, os.Status,
	).Scan(&os.ID, &os.CriadoEm, &os.AtualizadoEm)
}

func (r *OrdemServicoRepository) Atualizar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		UPDATE ordens_servico
		SET numero_ordem = $1, cliente_id = $2, defeito = $3, servicos = $4, assinatura = $5, status = $6, atualizado_em = NOW()
		WHERE id = $7
		RETURNING atualizado_em
	`
	return r.Pool.QueryRow(
		ctx, query,
		os.NumeroOrdem, os.ClienteID, os.Defeito, os.Servicos, os.Assinatura, os.Status, os.ID,
	).Scan(&os.AtualizadoEm)
}

func (r *OrdemServicoRepository) AtualizarTx(ctx context.Context, tx pgx.Tx, os *models.OrdemServico) error {
	query := `
		UPDATE ordens_servico
		SET numero_ordem = $1, cliente_id = $2, defeito = $3, servicos = $4, assinatura = $5, status = $6, atualizado_em = NOW()
		WHERE id = $7
		RETURNING atualizado_em
	`
	return tx.QueryRow(
		ctx, query,
		os.NumeroOrdem, os.ClienteID, os.Defeito, os.Servicos, os.Assinatura, os.Status, os.ID,
	).Scan(&os.AtualizadoEm)
}

func (r *OrdemServicoRepository) Deletar(ctx context.Context, id int64) error {
	query := `DELETE FROM ordens_servico WHERE id = $1`
	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("ordem de serviço não encontrada")
	}
	return nil
}

func (r *OrdemServicoRepository) DeletarTx(ctx context.Context, tx pgx.Tx, id int64) error {
	query := `DELETE FROM ordens_servico WHERE id = $1`
	cmdTag, err := tx.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errors.New("ordem de serviço não encontrada")
	}
	return nil
}

// ============================================================================
// MÉTODOS DE LEITURA
// ============================================================================

func (r *OrdemServicoRepository) BuscarPorID(ctx context.Context, id int64) (*models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE id = $1
	`
	var os models.OrdemServico
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&os.ID, &os.NumeroOrdem, &os.ClienteID, &os.Defeito, &os.Servicos, &os.Assinatura, &os.Status, &os.CriadoEm, &os.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &os, nil
}

func (r *OrdemServicoRepository) BuscarPorIDTx(ctx context.Context, tx pgx.Tx, id int64) (*models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE id = $1
		FOR UPDATE
	`
	var os models.OrdemServico
	err := tx.QueryRow(ctx, query, id).Scan(
		&os.ID, &os.NumeroOrdem, &os.ClienteID, &os.Defeito, &os.Servicos, &os.Assinatura, &os.Status, &os.CriadoEm, &os.AtualizadoEm,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &os, nil
}

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

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) BuscarPorClienteID(ctx context.Context, clienteID int64, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE cliente_id = $1
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.Pool.Query(ctx, query, clienteID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) BuscarPorNomeCliente(ctx context.Context, nome string, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT os.id, os.numero_ordem, os.cliente_id, os.defeito, os.servicos, os.assinatura, os.status, os.criado_em, os.atualizado_em
		FROM ordens_servico os
		JOIN clientes c ON os.cliente_id = c.id
		WHERE c.nome ILIKE $1
		ORDER BY os.id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.Pool.Query(ctx, query, "%"+nome+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) BuscarPorEmailCliente(ctx context.Context, email string, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT os.id, os.numero_ordem, os.cliente_id, os.defeito, os.servicos, os.assinatura, os.status, os.criado_em, os.atualizado_em
		FROM ordens_servico os
		JOIN clientes c ON os.cliente_id = c.id
		WHERE c.email ILIKE $1
		ORDER BY os.id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.Pool.Query(ctx, query, "%"+email+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) BuscarPorNumeroCliente(ctx context.Context, numero string, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT os.id, os.numero_ordem, os.cliente_id, os.defeito, os.servicos, os.assinatura, os.status, os.criado_em, os.atualizado_em
		FROM ordens_servico os
		JOIN clientes c ON os.cliente_id = c.id
		WHERE c.telefone ILIKE $1
		ORDER BY os.id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.Pool.Query(ctx, query, "%"+numero+"%", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) ListarPorPeriodo(ctx context.Context, dataInicio, dataFim time.Time) ([]models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, defeito, servicos, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE criado_em >= $1 AND criado_em <= $2
		ORDER BY criado_em DESC
	`
	rows, err := r.Pool.Query(ctx, query, dataInicio, dataFim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) UltimosTresMeses(ctx context.Context) ([]models.OrdemServico, error) {
	agora := time.Now()
	tresMesesAtras := agora.AddDate(0, -3, 0)
	return r.ListarPorPeriodo(ctx, tresMesesAtras, agora)
}

func (r *OrdemServicoRepository) UltimosSeisMeses(ctx context.Context) ([]models.OrdemServico, error) {
	agora := time.Now()
	seisMesesAtras := agora.AddDate(0, -6, 0)
	return r.ListarPorPeriodo(ctx, seisMesesAtras, agora)
}

// Helper para mapeamento dos resultados
func (r *OrdemServicoRepository) scanRows(rows pgx.Rows) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	for rows.Next() {
		var os models.OrdemServico
		if err := rows.Scan(
			&os.ID, &os.NumeroOrdem, &os.ClienteID, &os.Defeito, &os.Servicos, &os.Assinatura, &os.Status, &os.CriadoEm, &os.AtualizadoEm,
		); err != nil {
			return nil, err
		}
		ordens = append(ordens, os)
	}
	return ordens, rows.Err()
}

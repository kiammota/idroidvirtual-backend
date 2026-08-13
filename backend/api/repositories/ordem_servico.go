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
	return &OrdemServicoRepository{
		Pool: pool,
	}
}

// ============================================================================
// MÉTODOS DE ESCRITA / MUTACIONAL
// ============================================================================

func (r *OrdemServicoRepository) Criar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		INSERT INTO ordens_servico (usuario_id, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, numero_ordem, criado_em, atualizado_em
	`
	return r.Pool.QueryRow(
		ctx, query,
		os.UsuarioID, os.ClienteID, os.NomeServico, os.Defeito, os.Servicos, os.Valor, os.Custo, os.Assinatura, os.Status,
	).Scan(&os.ID, &os.NumeroOrdem, &os.CriadoEm, &os.AtualizadoEm)
}

func (r *OrdemServicoRepository) CriarTx(ctx context.Context, tx pgx.Tx, os *models.OrdemServico) error {
	query := `
		INSERT INTO ordens_servico (usuario_id, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, numero_ordem, criado_em, atualizado_em
	`
	return tx.QueryRow(
		ctx, query,
		os.UsuarioID, os.ClienteID, os.NomeServico, os.Defeito, os.Servicos, os.Valor, os.Custo, os.Assinatura, os.Status,
	).Scan(&os.ID, &os.NumeroOrdem, &os.CriadoEm, &os.AtualizadoEm)
}

func (r *OrdemServicoRepository) Atualizar(ctx context.Context, os *models.OrdemServico) error {
	query := `
		UPDATE ordens_servico
		SET cliente_id = $1, nome_servico = $2, defeito = $3, servicos = $4, valor = $5, custo = $6, assinatura = $7, status = $8, atualizado_em = NOW()
		WHERE id = $9
		RETURNING atualizado_em
	`
	return r.Pool.QueryRow(
		ctx, query,
		os.ClienteID, os.NomeServico, os.Defeito, os.Servicos, os.Valor, os.Custo, os.Assinatura, os.Status, os.ID,
	).Scan(&os.AtualizadoEm)
}

func (r *OrdemServicoRepository) AtualizarTx(ctx context.Context, tx pgx.Tx, os *models.OrdemServico) error {
	query := `
		UPDATE ordens_servico
		SET cliente_id = $1, nome_servico = $2, defeito = $3, servicos = $4, valor = $5, custo = $6, assinatura = $7, status = $8, atualizado_em = NOW()
		WHERE id = $9
		RETURNING atualizado_em
	`
	return tx.QueryRow(
		ctx, query,
		os.ClienteID, os.NomeServico, os.Defeito, os.Servicos, os.Valor, os.Custo, os.Assinatura, os.Status, os.ID,
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
		SELECT
			os.id,
			os.numero_ordem,
			os.cliente_id,
			os.nome_servico,
			os.defeito,
			os.servicos,
			os.valor,
			os.custo,
			os.assinatura,
			os.status,
			os.criado_em,
			os.atualizado_em,
			c.telefone
		FROM ordens_servico os
		JOIN clientes c ON c.id = os.cliente_id
		WHERE os.id = $1
	`

	var os models.OrdemServico

	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&os.ID,
		&os.NumeroOrdem,
		&os.ClienteID,
		&os.NomeServico,
		&os.Defeito,
		&os.Servicos,
		&os.Valor,
		&os.Custo,
		&os.Assinatura,
		&os.Status,
		&os.CriadoEm,
		&os.AtualizadoEm,
		&os.ClienteTelefone,
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
		SELECT id, numero_ordem, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE id = $1
		FOR UPDATE
	`
	var os models.OrdemServico
	err := tx.QueryRow(ctx, query, id).Scan(
		&os.ID, &os.NumeroOrdem, &os.ClienteID, &os.NomeServico, &os.Defeito, &os.Servicos, &os.Valor, &os.Custo, &os.Assinatura, &os.Status, &os.CriadoEm, &os.AtualizadoEm,
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
		SELECT id, numero_ordem, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status, criado_em, atualizado_em
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
		SELECT id, numero_ordem, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status, criado_em, atualizado_em
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
		SELECT os.id, os.numero_ordem, os.cliente_id, os.nome_servico, os.defeito, os.servicos, os.valor, os.custo, os.assinatura, os.status, os.criado_em, os.atualizado_em
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
		SELECT os.id, os.numero_ordem, os.cliente_id, os.nome_servico, os.defeito, os.servicos, os.valor, os.custo, os.assinatura, os.status, os.criado_em, os.atualizado_em
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

func (r *OrdemServicoRepository) ListarPorStatus(ctx context.Context, status string, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT id, numero_ordem, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status, criado_em, atualizado_em
		FROM ordens_servico
		WHERE ($1 = '' OR status = $1)
		ORDER BY id DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.Pool.Query(ctx, query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *OrdemServicoRepository) BuscarPorNumeroCliente(ctx context.Context, numero string, limit, offset int) ([]models.OrdemServico, error) {
	query := `
		SELECT os.id, os.numero_ordem, os.cliente_id, os.nome_servico, os.defeito, os.servicos, os.valor, os.custo, os.assinatura, os.status, os.criado_em, os.atualizado_em
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
		SELECT id, numero_ordem, cliente_id, nome_servico, defeito, servicos, valor, custo, assinatura, status, criado_em, atualizado_em
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

func (r *OrdemServicoRepository) ContarEmAndamentoTx(ctx context.Context, tx pgx.Tx) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM ordens_servico WHERE status = 'EM_ANDAMENTO'`
	err := tx.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
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

// ============================================================================
// HELPER DE MAPEAMENTO DOS RESULTADOS
// ============================================================================

func (r *OrdemServicoRepository) scanRows(rows pgx.Rows) ([]models.OrdemServico, error) {
	var ordens []models.OrdemServico
	for rows.Next() {
		var os models.OrdemServico
		err := rows.Scan(
			&os.ID,
			&os.NumeroOrdem,
			&os.ClienteID,
			&os.NomeServico,
			&os.Defeito,
			&os.Servicos,
			&os.Valor,
			&os.Custo,
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

package repositories

import (
	"context"
	"errors"
	"fmt"

	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ClienteRepository struct {
	Pool *pgxpool.Pool
}

func NewClienteRepository(pool *pgxpool.Pool) ClienteRepository {
	return ClienteRepository{Pool: pool}
}

// Criar insere um novo cliente
func (r *ClienteRepository) Criar(ctx context.Context, c *models.Cliente) error {

	query := `
		INSERT INTO clientes (
			nome,
			telefone,
			endereco,
			email
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING id, criado_em, atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		c.Nome,
		c.Telefone,
		c.Endereco,
		c.Email,
	).Scan(
		&c.ID,
		&c.CriadoEm,
		&c.AtualizadoEm,
	)

	if err != nil {
		return err
	}

	return nil
}

// BuscarPorTelefoneTx busca um cliente exatamente pelo número de telefone dentro de uma transação
func (r *ClienteRepository) BuscarPorTelefoneTx(
	ctx context.Context,
	tx pgx.Tx,
	telefone string,
) (*models.Cliente, error) {

	query := `
		SELECT 
			id,
			nome,
			telefone,
			endereco,
			email,
			criado_em,
			atualizado_em
		FROM clientes
		WHERE telefone = $1
		FOR UPDATE
	`

	var c models.Cliente

	err := tx.QueryRow(
		ctx,
		query,
		telefone,
	).Scan(
		&c.ID,
		&c.Nome,
		&c.Telefone,
		&c.Endereco,
		&c.Email,
		&c.CriadoEm,
		&c.AtualizadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

// CriarTx insere um novo cliente dentro de uma transação e popula seu ID e datas
func (r *ClienteRepository) CriarTx(
	ctx context.Context,
	tx pgx.Tx,
	c *models.Cliente,
) error {

	query := `
		INSERT INTO clientes (
			nome,
			telefone,
			endereco,
			email
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING id, criado_em, atualizado_em
	`

	err := tx.QueryRow(
		ctx,
		query,
		c.Nome,
		c.Telefone,
		c.Endereco,
		c.Email,
	).Scan(
		&c.ID,
		&c.CriadoEm,
		&c.AtualizadoEm,
	)

	if err != nil {
		return err
	}

	return nil
}

// BuscarPorID retorna um cliente pelo ID
func (r *ClienteRepository) BuscarPorID(ctx context.Context, id int64) (*models.Cliente, error) {

	query := `
		SELECT 
			id,
			nome,
			telefone,
			endereco,
			email,
			criado_em,
			atualizado_em
		FROM clientes
		WHERE id = $1
	`

	var c models.Cliente

	err := r.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&c.ID,
		&c.Nome,
		&c.Telefone,
		&c.Endereco,
		&c.Email,
		&c.CriadoEm,
		&c.AtualizadoEm,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &c, nil
}

// Existe verifica se um cliente existe
func (r *ClienteRepository) Existe(ctx context.Context, id int64) (bool, error) {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM clientes
			WHERE id = $1
		)
	`

	var existe bool

	err := r.Pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(&existe)

	return existe, err
}

// Listar retorna clientes paginados
func (r *ClienteRepository) Listar(ctx context.Context, limit, offset int) ([]models.Cliente, error) {

	query := `
		SELECT
			id,
			nome,
			telefone,
			endereco,
			email,
			criado_em,
			atualizado_em
		FROM clientes
		ORDER BY nome ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var clientes []models.Cliente

	for rows.Next() {

		var c models.Cliente

		err := rows.Scan(
			&c.ID,
			&c.Nome,
			&c.Telefone,
			&c.Endereco,
			&c.Email,
			&c.CriadoEm,
			&c.AtualizadoEm,
		)

		if err != nil {
			return nil, err
		}

		clientes = append(clientes, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return clientes, nil
}

// BuscarPorNome pesquisa clientes pelo nome
func (r *ClienteRepository) BuscarPorNome(
	ctx context.Context,
	nome string,
	limit,
	offset int,
) ([]models.Cliente, error) {

	query := `
		SELECT
			id,
			nome,
			telefone,
			endereco,
			email,
			criado_em,
			atualizado_em
		FROM clientes
		WHERE nome ILIKE $1
		ORDER BY nome ASC
		LIMIT $2 OFFSET $3
	`

	termoBusca := fmt.Sprintf("%%%s%%", nome)

	rows, err := r.Pool.Query(
		ctx,
		query,
		termoBusca,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var clientes []models.Cliente

	for rows.Next() {

		var c models.Cliente

		err := rows.Scan(
			&c.ID,
			&c.Nome,
			&c.Telefone,
			&c.Endereco,
			&c.Email,
			&c.CriadoEm,
			&c.AtualizadoEm,
		)

		if err != nil {
			return nil, err
		}

		clientes = append(clientes, c)
	}

	return clientes, nil
}

// Atualizar altera os dados do cliente
func (r *ClienteRepository) Atualizar(ctx context.Context, c *models.Cliente) error {

	query := `
		UPDATE clientes
		SET
			nome = $1,
			telefone = $2,
			endereco = $3,
			email = $4,
			atualizado_em = NOW()
		WHERE id = $5
		RETURNING atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		c.Nome,
		c.Telefone,
		c.Endereco,
		c.Email,
		c.ID,
	).Scan(
		&c.AtualizadoEm,
	)

	if err != nil {
		return err
	}

	return nil
}

// Deletar remove um cliente
func (r *ClienteRepository) Deletar(ctx context.Context, id int64) error {

	query := `
		DELETE FROM clientes
		WHERE id = $1
	`

	cmdTag, err := r.Pool.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum cliente encontrado para deletar")
	}

	return nil
}

func (r *ClienteRepository) BuscarPorIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
) (*models.Cliente, error) {

	query := `
		SELECT 
			id,
			nome,
			telefone,
			endereco,
			email,
			criado_em,
			atualizado_em
		FROM clientes
		WHERE id = $1
		FOR UPDATE
	`

	var c models.Cliente

	err := tx.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&c.ID,
		&c.Nome,
		&c.Telefone,
		&c.Endereco,
		&c.Email,
		&c.CriadoEm,
		&c.AtualizadoEm,
	)

	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &c, nil
}

// BuscarPorTelefone busca clientes por número de telefone
func (r *ClienteRepository) BuscarPorTelefone(ctx context.Context, telefone string) ([]models.Cliente, error) {
	query := `
		SELECT id, nome, telefone, endereco, email, criado_em, atualizado_em
		FROM clientes
		WHERE telefone ILIKE $1
		ORDER BY nome ASC
	`
	termoBusca := fmt.Sprintf("%%%s%%", telefone)
	rows, err := r.Pool.Query(ctx, query, termoBusca)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Telefone, &c.Endereco, &c.Email, &c.CriadoEm, &c.AtualizadoEm); err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

// BuscarPorEndereco busca clientes por endereço
func (r *ClienteRepository) BuscarPorEndereco(ctx context.Context, endereco string) ([]models.Cliente, error) {
	query := `
		SELECT id, nome, telefone, endereco, email, criado_em, atualizado_em
		FROM clientes
		WHERE endereco ILIKE $1
		ORDER BY nome ASC
	`
	termoBusca := fmt.Sprintf("%%%s%%", endereco)
	rows, err := r.Pool.Query(ctx, query, termoBusca)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Telefone, &c.Endereco, &c.Email, &c.CriadoEm, &c.AtualizadoEm); err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

// BuscarPorEmail busca clientes pelo e-mail
func (r *ClienteRepository) BuscarPorEmail(ctx context.Context, email string) ([]models.Cliente, error) {
	query := `
		SELECT id, nome, telefone, endereco, email, criado_em, atualizado_em
		FROM clientes
		WHERE email ILIKE $1
		ORDER BY nome ASC
	`
	termoBusca := fmt.Sprintf("%%%s%%", email)
	rows, err := r.Pool.Query(ctx, query, termoBusca)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clientes []models.Cliente
	for rows.Next() {
		var c models.Cliente
		if err := rows.Scan(&c.ID, &c.Nome, &c.Telefone, &c.Endereco, &c.Email, &c.CriadoEm, &c.AtualizadoEm); err != nil {
			return nil, err
		}
		clientes = append(clientes, c)
	}
	return clientes, nil
}

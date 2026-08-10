package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/models"
)

type ProdutoRepository struct {
	Pool *pgxpool.Pool
}

func NewProdutoRepository(pool *pgxpool.Pool) ProdutoRepository {
	return ProdutoRepository{Pool: pool}
}

// Criar insere um novo produto no banco de dados e retorna o ID e as datas geradas
func (r *ProdutoRepository) Criar(ctx context.Context, p *models.Produto) error {
	query := `
		INSERT INTO produtos (
			nome, preco, categoria, estoque
		) VALUES (
			$1, $2, $3, $4
		)
		RETURNING id, criado_em, atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		p.Nome,
		p.Preco,
		p.Categoria,
		p.Estoque,
	).Scan(&p.ID, &p.CriadoEm, &p.AtualizadoEm)

	if err != nil {
		return err
	}

	return nil
}

// BuscarPorID retorna um produto pelo seu ID
func (r *ProdutoRepository) BuscarPorID(ctx context.Context, id int64) (*models.Produto, error) {
	query := `
		SELECT id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE id = $1
	`

	var p models.Produto
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Nome,
		&p.Preco,
		&p.Categoria,
		&p.Estoque,
		&p.CriadoEm,
		&p.AtualizadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Retorna nil se o produto não for encontrado
		}
		return nil, err
	}

	return &p, nil
}

// BuscarPorIDTx busca por ID dentro de uma transação com trava (FOR UPDATE)
func (r *ProdutoRepository) BuscarPorIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
) (*models.Produto, error) {

	query := `
		SELECT id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE id = $1
		FOR UPDATE
	`

	var p models.Produto

	err := tx.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.Nome,
		&p.Preco,
		&p.Categoria,
		&p.Estoque,
		&p.CriadoEm,
		&p.AtualizadoEm,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &p, nil
}

// Listar retorna os produtos com suporte a paginação
func (r *ProdutoRepository) Listar(ctx context.Context, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		ORDER BY nome ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.Nome,
			&p.Preco,
			&p.Categoria,
			&p.Estoque,
			&p.CriadoEm,
			&p.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return produtos, nil
}

// BuscarPorNome faz uma busca parcial por nome
func (r *ProdutoRepository) BuscarPorNome(ctx context.Context, nome string, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE nome ILIKE $1
		ORDER BY nome ASC
		LIMIT $2 OFFSET $3
	`

	termoBusca := fmt.Sprintf("%%%s%%", nome)
	rows, err := r.Pool.Query(ctx, query, termoBusca, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.Nome,
			&p.Preco,
			&p.Categoria,
			&p.Estoque,
			&p.CriadoEm,
			&p.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, p)
	}

	return produtos, nil
}

// BuscarPorCategoria faz busca por categoria (case-insensitive)
func (r *ProdutoRepository) BuscarPorCategoria(ctx context.Context, categoria string, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE categoria ILIKE $1
		ORDER BY nome ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, categoria, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.Nome,
			&p.Preco,
			&p.Categoria,
			&p.Estoque,
			&p.CriadoEm,
			&p.AtualizadoEm,
		)
		if err != nil {
			return nil, err
		}
		produtos = append(produtos, p)
	}

	return produtos, nil
}

// Atualizar altera nome, preço, categoria e estoque do produto
func (r *ProdutoRepository) Atualizar(ctx context.Context, p *models.Produto) error {
	query := `
		UPDATE produtos
		SET nome = $1,
		    preco = $2,
		    categoria = $3,
		    estoque = $4,
		    atualizado_em = NOW()
		WHERE id = $5
		RETURNING atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		p.Nome,
		p.Preco,
		p.Categoria,
		p.Estoque,
		p.ID,
	).Scan(&p.AtualizadoEm)

	if err != nil {
		return err
	}

	return nil
}

// AumentarEstoque incrementa a quantidade do produto no estoque
func (r *ProdutoRepository) AumentarEstoque(ctx context.Context, id int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser adicionada deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque + $1,
		    atualizado_em = NOW()
		WHERE id = $2
	`

	cmdTag, err := r.Pool.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para aumentar o estoque")
	}

	return nil
}

func (r *ProdutoRepository) AlterarEstoque(
	ctx context.Context,
	id int64,
	quantidade int,
) error {
	query := `
		UPDATE produtos
		SET estoque = estoque + $1,
		    atualizado_em = NOW()
		WHERE id = $2
	`

	cmdTag, err := r.Pool.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado")
	}

	return nil
}

// AumentarEstoqueTx incrementa o estoque dentro de uma transação
func (r *ProdutoRepository) AumentarEstoqueTx(ctx context.Context, tx pgx.Tx, id int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser adicionada deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque + $1,
		    atualizado_em = NOW()
		WHERE id = $2
	`

	cmdTag, err := tx.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para aumentar o estoque")
	}

	return nil
}

// DiminuirEstoque decremanta o estoque garantindo valor não negativo
func (r *ProdutoRepository) DiminuirEstoque(ctx context.Context, id int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser removida deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque - $1,
		    atualizado_em = NOW()
		WHERE id = $2 AND estoque >= $1
	`

	cmdTag, err := r.Pool.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		produto, errBuscar := r.BuscarPorID(ctx, id)
		if errBuscar != nil || produto == nil {
			return errors.New("nenhum produto encontrado para diminuir o estoque")
		}
		return errors.New("estoque insuficiente para realizar a baixa")
	}

	return nil
}

// DiminuirEstoqueTx diminui o estoque dentro de uma transação
func (r *ProdutoRepository) DiminuirEstoqueTx(ctx context.Context, tx pgx.Tx, id int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser removida deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque - $1,
		    atualizado_em = NOW()
		WHERE id = $2 AND estoque >= $1
	`

	cmdTag, err := tx.Exec(ctx, query, quantidade, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		var existe bool
		errExiste := tx.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1)", id).Scan(&existe)
		if errExiste != nil || !existe {
			return errors.New("nenhum produto encontrado para diminuir o estoque")
		}
		return errors.New("estoque insuficiente para realizar a baixa")
	}

	return nil
}

// Deletar remove um produto pelo ID
func (r *ProdutoRepository) Deletar(ctx context.Context, id int64) error {
	query := `DELETE FROM produtos WHERE id = $1`

	cmdTag, err := r.Pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para deletar")
	}

	return nil
}

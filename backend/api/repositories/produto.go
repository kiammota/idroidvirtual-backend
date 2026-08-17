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

// Criar insere um novo produto
func (r *ProdutoRepository) Criar(ctx context.Context, p *models.Produto) error {
	query := `
		INSERT INTO produtos (
			empresa_id, nome, preco, categoria, estoque
		) VALUES (
			$1, $2, $3, $4, $5
		)
		RETURNING id, criado_em, atualizado_em
	`

	err := r.Pool.QueryRow(
		ctx,
		query,
		p.EmpresaID,
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

// BuscarPorID retorna um produto pelo ID (dentro da empresa)
func (r *ProdutoRepository) BuscarPorID(ctx context.Context, id int) (*models.Produto, error) {
	query := `
		SELECT id, empresa_id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE id = $1
	`

	var p models.Produto
	err := r.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.EmpresaID,
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

// BuscarPorIDTx busca por ID dentro de uma transação com trava (FOR UPDATE)
func (r *ProdutoRepository) BuscarPorIDTx(
	ctx context.Context,
	tx pgx.Tx,
	id int64,
) (*models.Produto, error) {

	query := `
		SELECT id, empresa_id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE id = $1 
		FOR UPDATE
	`

	var p models.Produto

	err := tx.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.EmpresaID,
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

// Listar retorna os produtos da empresa com paginação
func (r *ProdutoRepository) Listar(ctx context.Context, empresaID int64, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, empresa_id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE empresa_id = $1
		ORDER BY nome ASC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.EmpresaID,
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

// BuscarPorNome faz busca parcial por nome dentro da empresa
func (r *ProdutoRepository) BuscarPorNome(ctx context.Context, empresaID int64, nome string, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, empresa_id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE empresa_id = $1 AND nome ILIKE $2
		ORDER BY nome ASC
		LIMIT $3 OFFSET $4
	`

	termoBusca := fmt.Sprintf("%%%s%%", nome)
	rows, err := r.Pool.Query(ctx, query, empresaID, termoBusca, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.EmpresaID,
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

// BuscarPorCategoria faz busca por categoria dentro da empresa
func (r *ProdutoRepository) BuscarPorCategoria(ctx context.Context, empresaID int64, categoria string, limit, offset int) ([]models.Produto, error) {
	query := `
		SELECT id, empresa_id, nome, preco, categoria, estoque, criado_em, atualizado_em
		FROM produtos
		WHERE empresa_id = $1 AND categoria ILIKE $2
		ORDER BY nome ASC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.Pool.Query(ctx, query, empresaID, categoria, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var produtos []models.Produto
	for rows.Next() {
		var p models.Produto
		err := rows.Scan(
			&p.ID,
			&p.EmpresaID,
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

// Atualizar altera nome, preço, categoria e estoque apenas pelo ID do produto
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

// AlterarEstoqueTx define o estoque de um produto para uma quantidade específica dentro de uma transação, filtrando pela empresa
func (r *ProdutoRepository) AlterarEstoqueTx(ctx context.Context, tx pgx.Tx, empresaID, produtoID int64, novaQuantidade int) error {
	if novaQuantidade < 0 {
		return errors.New("a quantidade em estoque não pode ser negativa")
	}

	query := `
		UPDATE produtos
		SET estoque = $1,
			atualizado_em = NOW()
		WHERE id = $2 AND empresa_id = $3
	`

	cmdTag, err := tx.Exec(ctx, query, novaQuantidade, produtoID, empresaID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para alterar o estoque ou o produto não pertence a esta empresa")
	}

	return nil
}

// AlterarEstoque define o estoque de um produto para uma quantidade específica diretamente pelo pool, filtrando pela empresa
func (r *ProdutoRepository) AlterarEstoque(ctx context.Context, produtoID int64, novaQuantidade int) error {
	if novaQuantidade < 0 {
		return errors.New("a quantidade em estoque não pode ser negativa")
	}

	query := `
		UPDATE produtos
		SET estoque = $1,
			atualizado_em = NOW()
		WHERE id = $2
	`

	cmdTag, err := r.Pool.Exec(ctx, query, novaQuantidade, produtoID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para alterar o estoque ou o produto não pertence a esta empresa")
	}

	return nil
}

// AumentarEstoque incrementa a quantidade do produto apenas pelo ID
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

// DiminuirEstoque decrementa o estoque garantindo valor não negativo apenas pelo ID
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
		var existe bool
		errExiste := r.Pool.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM produtos WHERE id = $1)", id).Scan(&existe)
		if errExiste != nil || !existe {
			return errors.New("nenhum produto encontrado para diminuir o estoque")
		}
		return errors.New("estoque insuficiente para realizar a baixa")
	}

	return nil
}

// DiminuirEstoqueTx diminui o estoque dentro de uma transação apenas pelo ID
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

func (r *ProdutoRepository) AumentarEstoqueTx(ctx context.Context, tx pgx.Tx, produtoID int64, quantidade int) error {
	if quantidade <= 0 {
		return errors.New("a quantidade a ser adicionada deve ser maior que zero")
	}

	query := `
		UPDATE produtos
		SET estoque = estoque + $1,
			atualizado_em = NOW()
		WHERE id = $2 
	`

	cmdTag, err := tx.Exec(ctx, query, quantidade, produtoID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("nenhum produto encontrado para aumentar o estoque ou o produto não pertence a esta empresa")
	}

	return nil
}

// Deletar remove um produto do banco de dados apenas pelo ID
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

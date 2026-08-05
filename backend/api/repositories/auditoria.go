package repositories

import (
	"context"
	"encoding/json"
	"time"

	"backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditoriaRepository struct {
	db *pgxpool.Pool
}

func NewAuditoriaRepository(db *pgxpool.Pool) *AuditoriaRepository {
	return &AuditoriaRepository{db: db}
}

func (r *AuditoriaRepository) Criar(ctx context.Context, audit models.Auditoria) error {
	var valoresAntigos, valoresNovos []byte
	var err error

	if audit.ValoresAntigos != nil {
		valoresAntigos, err = json.Marshal(audit.ValoresAntigos)
		if err != nil {
			return err
		}
	}

	if audit.ValoresNovos != nil {
		valoresNovos, err = json.Marshal(audit.ValoresNovos)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO auditorias (
			usuario_id, acao, tabela, registro_id, valores_antigos, valores_novos, horario
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err = r.db.Exec(ctx, query,
		audit.UsuarioID,
		audit.Acao,
		audit.Tabela,
		audit.RegistroID,
		valoresAntigos,
		valoresNovos,
	)
	return err
}

func (r *AuditoriaRepository) CriarTx(ctx context.Context, tx pgx.Tx, audit models.Auditoria) error {
	var valoresAntigos, valoresNovos []byte
	var err error

	if audit.ValoresAntigos != nil {
		valoresAntigos, err = json.Marshal(audit.ValoresAntigos)
		if err != nil {
			return err
		}
	}

	if audit.ValoresNovos != nil {
		valoresNovos, err = json.Marshal(audit.ValoresNovos)
		if err != nil {
			return err
		}
	}

	query := `
		INSERT INTO auditorias (
			usuario_id, acao, tabela, registro_id, valores_antigos, valores_novos, horario
		) VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`

	_, err = tx.Exec(ctx, query,
		audit.UsuarioID,
		audit.Acao,
		audit.Tabela,
		audit.RegistroID,
		valoresAntigos,
		valoresNovos,
	)
	return err
}

func (r *AuditoriaRepository) BuscarPorPeriodo(ctx context.Context, inicio, fim time.Time, limit, offset int) ([]models.Auditoria, error) {
	query := `
		SELECT id, usuario_id, acao, tabela, registro_id, valores_antigos, valores_novos, horario
		FROM auditorias
		WHERE horario BETWEEN $1 AND $2
		ORDER BY horario DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, inicio, fim, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var lista []models.Auditoria
	for rows.Next() {
		var a models.Auditoria
		var antigos, novos []byte

		err := rows.Scan(
			&a.ID,
			&a.UsuarioID,
			&a.Acao,
			&a.Tabela,
			&a.RegistroID,
			&antigos,
			&novos,
			&a.Horario,
		)
		if err != nil {
			return nil, err
		}

		if len(antigos) > 0 {
			_ = json.Unmarshal(antigos, &a.ValoresAntigos)
		}
		if len(novos) > 0 {
			_ = json.Unmarshal(novos, &a.ValoresNovos)
		}

		lista = append(lista, a)
	}

	return lista, nil
}

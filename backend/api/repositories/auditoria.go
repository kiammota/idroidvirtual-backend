package repositories

import (
	"context"
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
	query := `
		INSERT INTO auditorias (usuario_id, acao, horario) 
		VALUES ($1, $2, NOW())
	`
	_, err := r.db.Exec(ctx, query, audit.UsuarioID, audit.Acao)
	return err
}

func (r *AuditoriaRepository) CriarTx(ctx context.Context, tx pgx.Tx, audit models.Auditoria) error {
	query := `
		INSERT INTO auditorias (usuario_id, acao, horario) 
		VALUES ($1, $2, NOW())
	`
	_, err := tx.Exec(ctx, query, audit.UsuarioID, audit.Acao)
	return err
}
func (r *AuditoriaRepository) BuscarPorPeriodo(ctx context.Context, inicio, fim time.Time, limit, offset int) ([]models.Auditoria, error) {
	query := `
		SELECT 
			a.id, 
			a.usuario_id, 
			COALESCE(u.nome, 'Usuário Desconhecido') AS nome_usuario,
			a.acao, 
			a.horario
		FROM auditorias a
		LEFT JOIN usuarios u ON u.id = a.usuario_id
		WHERE a.horario BETWEEN $1 AND $2
		ORDER BY a.horario DESC
		LIMIT $3 OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, inicio, fim, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Carrega a localização de São Paulo para formatar o retorno do Go
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.FixedZone("UTC-3", -3*60*60)
	}

	lista := make([]models.Auditoria, 0)
	for rows.Next() {
		var a models.Auditoria
		err := rows.Scan(
			&a.ID,
			&a.UsuarioID,
			&a.NomeUsuario,
			&a.Acao,
			&a.Horario,
		)
		if err != nil {
			return nil, err
		}

		// Converte a data lida do banco para a timezone de São Paulo
		a.Horario = a.Horario.In(loc)

		lista = append(lista, a)
	}

	return lista, nil
}

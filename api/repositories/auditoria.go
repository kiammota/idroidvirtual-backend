package repositories

import (
	"backend/models"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuditoriaRepository struct {
	DB *pgxpool.Pool
}

func NewAuditoriaRepository(
	db *pgxpool.Pool,
) *AuditoriaRepository {

	return &AuditoriaRepository{
		DB: db,
	}
}

func (r *AuditoriaRepository) Criar(
	ctx context.Context,
	auditoria models.Auditoria,
) error {

	valoresAntigos, err := json.Marshal(
		auditoria.ValoresAntigos,
	)

	if err != nil {
		return err
	}

	valoresNovos, err := json.Marshal(
		auditoria.ValoresNovos,
	)

	if err != nil {
		return err
	}

	query := `
		INSERT INTO auditorias (
			usuario_id,
			acao,
			tabela,
			registro_id,
			valores_antigos,
			valores_novos,
			horario
		)
		VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6,
			NOW()
		)
	`

	_, err = r.DB.Exec(
		ctx,
		query,
		auditoria.UsuarioID,
		auditoria.Acao,
		auditoria.Tabela,
		auditoria.RegistroID,
		valoresAntigos,
		valoresNovos,
	)

	return err
}

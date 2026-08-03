package main

import (
	"backend/models"
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
)

func CriarAuditoria(
	ctx context.Context,
	db *pgx.Conn,
	auditoria models.Auditoria,
) error {
	valoresAntigos, err := json.Marshal(auditoria.ValoresAntigos)
	if err != nil {
		return err
	}

	valoresNovos, err := json.Marshal(auditoria.ValoresNovos)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO auditoria (
			usuario_id,
			acao,
			tabela,
			registro_id,
			valores_antigos,
			valores_novos
		)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err = db.Exec(
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

package models

import "time"

type Auditoria struct {
	ID          int64     `json:"id"`
	EmpresaID   int64     `json:"empresa_id"`
	UsuarioID   int64     `json:"usuario_id"`
	NomeUsuario string    `json:"nome_usuario"`
	Acao        string    `json:"acao"`
	Horario     time.Time `json:"horario"`
}

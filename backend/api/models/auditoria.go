package models

import "time"

type Auditoria struct {
	ID             int64                  `json:"id"`
	UsuarioID      int64                  `json:"usuario_id"`
	Acao           string                 `json:"acao"`
	Tabela         string                 `json:"tabela"`
	RegistroID     int64                  `json:"registro_id"`
	ValoresAntigos map[string]interface{} `json:"valores_antigos"`
	ValoresNovos   map[string]interface{} `json:"valores_novos"`
	Horario        time.Time              `json:"horario"`
}

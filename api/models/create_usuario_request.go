package models

type CreateUsuarioRequest struct {
	EmpresaID int64  `json:"empresa_id"`
	Email     string `json:"email" binding:"required"`
	Nome      string `json:"nome" binding:"required"`
	Senha     string `json:"senha" binding:"required"`
}

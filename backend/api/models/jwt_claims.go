package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	EmpresaID int64  `json:"empresa_id"`
	UsuarioID int64  `json:"usuario_id"`
	Email     string `json:"email"`

	jwt.RegisteredClaims
}

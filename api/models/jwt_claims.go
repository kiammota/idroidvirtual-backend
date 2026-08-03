package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UsuarioID int64  `json:"usuario_id"`
	Email     string `json:"email"`

	jwt.RegisteredClaims
}

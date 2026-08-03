package services

import (
	"backend/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	SecretKey string
}

func NewTokenService(secret string) *TokenService {
	return &TokenService{
		SecretKey: secret,
	}
}

// Gera um JWT
func (s *TokenService) GerarToken(usuarioID int64, email string) (string, error) {

	claims := models.JWTClaims{
		UsuarioID: usuarioID,
		Email:     email,

		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),

			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		&claims,
	)

	return token.SignedString(
		[]byte(s.SecretKey),
	)
}

// Valida e retorna os dados do usuário
func (s *TokenService) ValidarToken(tokenString string) (*models.JWTClaims, error) {

	claims := &models.JWTClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,

		func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("método de assinatura inválido")
			}

			return []byte(s.SecretKey), nil
		},
	)

	if err != nil || !token.Valid {
		return nil, errors.New("token inválido")
	}

	return claims, nil
}

package services

import (
	"errors"
	"time"

	"backend/models"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	SecretKey             string
	AccessTokenExpiration time.Duration
}

func NewTokenService(secret string, accessTokenExpiration time.Duration) *TokenService {
	return &TokenService{
		SecretKey:             secret,
		AccessTokenExpiration: accessTokenExpiration,
	}
}

// GerarToken gera um JWT contendo ID e Email do usuário.
func (s *TokenService) GerarToken(empresaID, usuarioID int64, email string) (string, error) {
	now := time.Now()

	claims := models.JWTClaims{
		EmpresaID: empresaID,
		UsuarioID: usuarioID,
		Email:     email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.SecretKey))
}

// ValidarToken valida a assinatura/expiração e retorna as claims já tipadas.
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
		return nil, errors.New("token inválido ou expirado")
	}

	return claims, nil
}

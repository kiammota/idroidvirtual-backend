package helpers

import (
	"backend/models"
	"errors"

	"github.com/gin-gonic/gin"
)

// GetClaims extrai as claims completas do contexto
func GetClaims(ctx *gin.Context) (*models.JWTClaims, error) {
	value, exists := ctx.Get("claims")
	if !exists {
		return nil, errors.New("claims não encontradas no contexto")
	}

	claims, ok := value.(*models.JWTClaims)
	if !ok {
		return nil, errors.New("formato de claims inválido")
	}

	return claims, nil
}

// GetEmpresaID retorna apenas o empresa_id
func GetEmpresaID(ctx *gin.Context) (int64, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return 0, err
	}
	return claims.EmpresaID, nil
}

// GetUserID retorna apenas o user_id
func GetUserID(ctx *gin.Context) (int64, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return 0, err
	}
	return claims.UsuarioID, nil
}

// GetUserEmail retorna o email do usuário
func GetUserEmail(ctx *gin.Context) (string, error) {
	claims, err := GetClaims(ctx)
	if err != nil {
		return "", err
	}
	return claims.Email, nil
}

// GetToken retorna o token JWT bruto (string)
func GetToken(ctx *gin.Context) (string, error) {
	value, exists := ctx.Get("token")
	if !exists {
		return "", errors.New("token não encontrado no contexto")
	}

	token, ok := value.(string)
	if !ok {
		return "", errors.New("token em formato inválido")
	}

	return token, nil
}

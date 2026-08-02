package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// GerarHashSenha cria um hash seguro da senha
func GerarHashSenha(senha string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerificarSenha compara a senha com o hash
func VerificarSenha(hash, senha string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(senha)) == nil
}

// GerarToken cria um JWT token válido por 24 horas
func GerarToken(usuario_id int64, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usuario_id": usuario_id,
		"email":      email,
		"exp":        time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(chave_jwt))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// VerificarToken valida e extrai claims do JWT
func VerificarToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(chave_jwt), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims, nil
}

// Middleware de autenticação
func Autenticado() gin.HandlerFunc {
	return func(c *gin.Context) {
		cabecalho_auth := c.GetHeader("Authorization")
		if cabecalho_auth == "" {
			c.JSON(401, gin.H{"erro": "autorização necessária"})
			c.Abort()
			return
		}

		var token_string string
		if len(cabecalho_auth) > 7 && cabecalho_auth[:7] == "Bearer " {
			token_string = cabecalho_auth[7:]
		} else {
			c.JSON(401, gin.H{"erro": "formato de autorização inválido"})
			c.Abort()
			return
		}

		claims, err := VerificarToken(token_string)
		if err != nil {
			c.JSON(401, gin.H{"erro": "token inválido"})
			c.Abort()
			return
		}

		c.Set("usuario_id", int64(claims["usuario_id"].(float64)))
		c.Set("email", claims["email"].(string))

		c.Next()
	}
}

// ObterUsuarioIDDoContexto extrai o usuario_id do contexto
func ObterUsuarioIDDoContexto(c *gin.Context) int64 {
	usuario_id, _ := c.Get("usuario_id")
	return usuario_id.(int64)
}

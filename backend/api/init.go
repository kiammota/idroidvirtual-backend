package main

import (
	"backend/logger"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	Porta         string
	URLBancoDados string
	ChaveJWT      string
}

func init() {
	logger.Log("Inicializando aplicação...")

	err := godotenv.Load()
	if err != nil {
		logger.Warn("Arquivo .env padrão não encontrado, continuando...")
		return
	}

	logger.OK("Arquivo .env carregado")
}

func CarregarConfig() Config {
	logger.Log("Iniciando carregamento das configurações")

	err := godotenv.Load("../.env")
	if err != nil {
		logger.Warn("Não foi possível carregar ../.env, usando variáveis de ambiente existentes")
	} else {
		logger.OK("../.env carregado")
	}

	logger.Log("Lendo variáveis do banco de dados")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbName == "" {
		logger.Warn("Algumas configurações do banco estão vazias")
	}

	logger.Log("Montando URL do banco de dados")

	url := os.Getenv("DATABASE_URL")

	if url == "" {
		logger.Log("DATABASE_URL não encontrada, criando URL automaticamente")

		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			dbName,
		)

		logger.OK("URL do banco criada")
	} else {
		logger.OK("DATABASE_URL encontrada")
	}

	porta := os.Getenv("PORT")

	if porta == "" {
		logger.Warn("PORT não configurada")
	}

	jwt := os.Getenv("JWT_SECRET")

	if jwt == "" {
		logger.Warn("JWT_SECRET não configurada")
	}

	logger.OK("Configurações carregadas com sucesso")

	return Config{
		Porta:         porta,
		URLBancoDados: url,
		ChaveJWT:      jwt,
	}
}

func (c Config) ExibirConfiguracao() {
	logger.Log("Exibindo configurações da aplicação")

	logger.Log("===== CONFIGURAÇÃO =====")
	logger.Log(fmt.Sprintf("Porta: %s", c.Porta))
	logger.Log(fmt.Sprintf("Banco configurado: %t", c.URLBancoDados != ""))
	logger.Log(fmt.Sprintf("JWT configurado: %t", c.ChaveJWT != ""))
	logger.Log("========================")

	logger.OK("Configurações exibidas")
}

func ConectarBanco(url string) *pgxpool.Pool {

init:
	logger.Log("Iniciando conexão com banco de dados")

	bd, err := pgxpool.New(context.Background(), url)

	if err != nil {
		logger.Err("Falha ao criar pool de conexão com banco")
		panic(err)
	}

	logger.Log("Testando conexão com banco")

	if err := bd.Ping(context.Background()); err != nil {
		logger.Err("Banco de dados não respondeu ao ping")
		logger.Warn("Tentando de novo em 5 segundos...")
		time.Sleep(5 * time.Second)
		goto init
	}

	logger.OK("Banco de dados conectado com sucesso")

	return bd
}

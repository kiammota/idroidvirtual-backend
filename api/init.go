package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	Porta         string
	URLBancoDados string
	ChaveJWT      string
}

func init() {
	_ = godotenv.Load()
}

func CarregarConfig() Config {
	err := godotenv.Load("../.env")

	println("carregando .env...")
	if err != nil {
		println("erro ao pegar o .env: ", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	url := os.Getenv("DATABASE_URL")
	if url == "" {
		url = fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			user,
			password,
			host,
			port,
			dbName,
		)
	}

	return Config{
		Porta:         os.Getenv("PORT"),
		URLBancoDados: url,
		ChaveJWT:      os.Getenv("JWT_SECRET"),
	}
}

func (c Config) ExibirConfiguracao() {
	log.Println("===== CONFIGURAÇÃO =====")
	log.Printf("Porta: %s", c.Porta)
	log.Printf("Banco: %s", c.URLBancoDados)
	log.Printf("JWT configurado: %t", c.ChaveJWT != "")
	log.Println("========================")
}

func ConectarBanco(url string) *pgxpool.Pool {
	bd, err := pgxpool.New(context.Background(), url)
	if err != nil {
		log.Fatalf("Falha ao conectar ao banco: %v", err)
	}

	if err := bd.Ping(context.Background()); err != nil {
		log.Fatalf("Falha ao fazer ping no banco: %v", err)
	}

	log.Println("Banco de dados conectado")
	return bd
}

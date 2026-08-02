package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var bd *pgxpool.Pool
var chave_jwt string

func init() {
	godotenv.Load()
}

type Configuracao struct {
	Porta         string
	URLBancoDados string
	ChaveJWT      string
}

func Carregar() Configuracao {

	err := godotenv.Load("../.env")
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

	return Configuracao{
		Porta:         os.Getenv("PORT"),
		URLBancoDados: url,
		ChaveJWT:      os.Getenv("JWT_SECRET"),
	}
}

func (c Configuracao) ExibirConfiguracao() {
	log.Println("===== CONFIGURAÇÃO =====")
	log.Printf("Porta: %s", c.Porta)
	log.Printf("Banco: %s", c.URLBancoDados)
	log.Printf("JWT configurado: %t", c.ChaveJWT != "")
	log.Println("========================")
}

func main() {
	// Carregar configurações
	cfg := Carregar()
	cfg.ExibirConfiguracao()

	// Conexão com banco
	var err error
	bd, err = pgxpool.New(context.Background(), cfg.URLBancoDados)
	if err != nil {
		log.Fatalf("❌ Falha ao conectar ao banco: %v", err)
	}
	defer bd.Close()

	if err := bd.Ping(context.Background()); err != nil {
		log.Fatalf("❌ Falha ao fazer ping no banco: %v", err)
	}
	log.Println("✅ Banco de dados conectado")

	chave_jwt = cfg.ChaveJWT

	// Setup Gin
	roteador := gin.Default()

	// CORS
	roteador.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// ===== ROTAS =====

	// Health check
	roteador.GET("/saude", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Autenticação (público)
	roteador.POST("/autenticar", AutenticarHandler)
	roteador.POST("/registrar", RegistrarHandler)

	// Usuários (protegido)
	roteador.GET("/usuarios/eu", Autenticado(), ObterMeuUsuarioHandler)

	// Produtos (protegido)
	roteador.GET("/produtos", Autenticado(), ListarProdutosHandler)
	roteador.POST("/produtos", Autenticado(), CriarProdutoHandler)
	roteador.PUT("/produtos/:id", Autenticado(), AtualizarProdutoHandler)
	roteador.DELETE("/produtos/:id", Autenticado(), DeletarProdutoHandler)

	// Clientes (protegido)
	roteador.GET("/clientes", Autenticado(), ListarClientesHandler)
	roteador.POST("/clientes", Autenticado(), CriarClienteHandler)
	roteador.PUT("/clientes/:id", Autenticado(), AtualizarClienteHandler)
	roteador.DELETE("/clientes/:id", Autenticado(), DeletarClienteHandler)

	// Vendas (protegido)
	roteador.GET("/vendas", Autenticado(), ListarVendasHandler)
	roteador.POST("/vendas", Autenticado(), CriarVendaHandler)

	// Ordens de Serviço (protegido)
	roteador.GET("/ordens-servico", Autenticado(), ListarOrdensServicoHandler)
	roteador.POST("/ordens-servico", Autenticado(), CriarOrdemServicoHandler)
	roteador.PUT("/ordens-servico/:id/status", Autenticado(), AtualizarStatusOrdemHandler)

	// Dashboard (protegido)
	roteador.GET("/painel", Autenticado(), ObterPainelHandler)

	// Iniciar servidor
	log.Printf("🚀 Servidor rodando na porta %s", cfg.Porta)
	roteador.Run("0.0.0.0:" + cfg.Porta)
}

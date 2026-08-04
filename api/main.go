package main

import (
	"log"

	"backend/gateways/controllers"
	"backend/gateways/middlewares"
	"backend/gateways/routes"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Configurações e Conexão ao Banco
	cfg := CarregarConfig()
	cfg.ExibirConfiguracao()

	bd := ConectarBanco(cfg.URLBancoDados)
	defer bd.Close()

	// 2. Repositórios
	auditoriaRepo := repositories.NewAuditoriaRepository(bd)
	usuarioRepo := repositories.NewUsuarioRepository(bd)
	vendaRepo := repositories.NewVendaRepository(bd)
	produtoRepo := repositories.NewProdutoRepository(bd)

	// 3. Serviços
	usuarioService := services.NewUsuarioService(usuarioRepo)
	tokenService := services.NewTokenService(cfg.ChaveJWT)
	loginService := services.NewLoginService(usuarioRepo, tokenService)
	vendaService := services.NewVendaService(bd, vendaRepo, produtoRepo, *auditoriaRepo)

	// 4. Controllers
	usuarioController := controllers.NewUsuarioController(usuarioService, loginService)
	vendaController := controllers.NewVendaController(vendaService)
	authController := controllers.NewAuthController(loginService, auditoriaRepo)

	// 5. Servidor HTTP e Middlewares Globais
	engine := gin.Default()
	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.JsonOnlyMiddleware())
	engine.Use(middlewares.JsonNotEmptyMiddleware())
	// public.Use(middlewares.LogMiddleware())

	// 6. Rotas Públicas
	engine.GET("/", func(ctx *gin.Context) { ctx.JSON(200, "pong!") })
	engine.POST("login", func(ctx *gin.Context) { authController.UsuarioLogar(ctx) })

	public := engine.Group("")
	private := engine.Group("")

	engine.Static("/swagger", "./swagger")

	private.Use(middlewares.AuthJWT())
	routes.AuthRoutes(public, authController)
	routes.UsuarioRoutes(public, usuarioController)
	routes.VendaRoutes(private, vendaController)

	// 7. Rotas Privadas (Autenticadas)

	// 8. Inicialização do Servidor
	println("rodando com swagger!")
	log.Printf("Servidor rodando na porta %s", cfg.Porta)
	if err := engine.Run("0.0.0.0:" + cfg.Porta); err != nil {
		log.Fatalf("Falha ao subir servidor: %v", err)
	}
}

package main

import (
	"backend/gateways/controllers"
	"backend/gateways/middlewares"
	"backend/gateways/routes"
	"backend/logger"
	"backend/repositories"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func main() {
	logger.Log("Iniciando aplicação")

	// 1. Configurações e Conexão ao Banco
	logger.Log("Carregando configurações")

	cfg := CarregarConfig()
	cfg.ExibirConfiguracao()

	logger.Log("Conectando ao banco de dados")

	bd := ConectarBanco(cfg.URLBancoDados)
	defer bd.Close()

	logger.OK("Conexão com banco estabelecida")

	// 2. Repositórios
	logger.Log("Inicializando repositórios")

	auditoriaRepo := repositories.NewAuditoriaRepository(bd)
	usuarioRepo := repositories.NewUsuarioRepository(bd)
	vendaRepo := repositories.NewVendaRepository(bd)
	produtoRepo := repositories.NewProdutoRepository(bd)
	ordemServicoRepo := repositories.NewOrdemServicoRepository(bd)
	clienteRepo := repositories.NewClienteRepository(bd)

	logger.OK("Repositórios carregados")

	// 3. Serviços
	logger.Log("Inicializando serviços")

	usuarioService := services.NewUsuarioService(usuarioRepo)
	tokenService := services.NewTokenService(cfg.ChaveJWT)
	loginService := services.NewLoginService(usuarioRepo, tokenService)
	vendaService := services.NewVendaService(bd, vendaRepo, produtoRepo, *auditoriaRepo)
	ordensServicoService := services.NewOrdensServicoService(
		bd,
		ordemServicoRepo,
		&produtoRepo,
		&clienteRepo,
	)

	logger.OK("Serviços carregados")

	// 4. Controllers
	logger.Log("Inicializando controllers")

	usuarioController := controllers.NewUsuarioController(usuarioService, loginService)
	vendaController := controllers.NewVendaController(vendaService)
	authController := controllers.NewAuthController(loginService, auditoriaRepo)
	ordensServicoController := controllers.NewOrdemServicoController(ordensServicoService)
	clienteController := controllers.NewClientController(&clienteRepo)

	logger.OK("Controllers carregados")

	// 5. Servidor HTTP e Middlewares Globais
	logger.Log("Configurando servidor HTTP")

	engine := gin.Default()

	logger.Log("Registrando middlewares")

	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.JsonOnlyMiddleware())
	engine.Use(middlewares.JsonNotEmptyMiddleware())

	logger.OK("Middlewares registrados")

	// 6. Rotas Públicas
	logger.Log("Registrando rotas públicas")

	engine.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, "pong!")
	})

	engine.POST("/login", func(ctx *gin.Context) {
		authController.UsuarioLogar(ctx)
	})

	public := engine.Group("")
	private := engine.Group("")

	logger.Log("Registrando Swagger")

	engine.Static("/swagger", "./swagger")

	logger.Log("Configurando rotas privadas")

	private.Use(middlewares.AuthJWT(tokenService))

	routes.AuthRoutes(public, authController)
	routes.UsuarioRoutes(public, usuarioController)
	routes.VendaRoutes(private, vendaController)
	routes.OrdensServicoRoutes(private, ordensServicoController)
	routes.ClienteRoutes(private, clienteController)

	logger.OK("Rotas registradas com sucesso")

	logger.OK("Swagger disponível")
	logger.Log("Servidor iniciando na porta " + cfg.Porta)

	if err := engine.Run("0.0.0.0:" + cfg.Porta); err != nil {
		logger.Err("Falha ao iniciar servidor: " + err.Error())
		panic(err)
	}
}

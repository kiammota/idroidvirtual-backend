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

	empresaRepo := repositories.NewEmpresaRepository(bd)
	auditoriaRepo := repositories.NewAuditoriaRepository(bd)
	usuarioRepo := repositories.NewUsuarioRepository(bd)
	vendaRepo := repositories.NewVendaRepository(bd)
	produtoRepo := repositories.NewProdutoRepository(bd)
	ordemServicoRepo := repositories.NewOrdemServicoRepository(bd)
	clienteRepo := repositories.NewClienteRepository(bd)
	trocaRepo := repositories.NewTrocaRepository(bd)
	financeiroRepo := repositories.NewFinanceiroRepository(bd)
	despesaRepo := repositories.NewDespesaRepository(bd)

	logger.OK("Repositórios carregados")

	// 3. Serviços
	logger.Log("Inicializando serviços")
	auditoriaService := services.NewAuditoriaService(auditoriaRepo, empresaRepo)
	usuarioService := services.NewUsuarioService(usuarioRepo)
	tokenService := services.NewTokenService(cfg.ChaveJWT)
	loginService := services.NewLoginService(usuarioRepo, tokenService)
	vendaService := services.NewVendaService(bd, vendaRepo, produtoRepo, *usuarioRepo, clienteRepo, *auditoriaService)
	ordensServicoService := services.NewOrdensServicoService(bd, ordemServicoRepo, &clienteRepo, auditoriaService)
	trocaService := services.NewTrocaService(trocaRepo)
	financeiroService := services.NewFinanceiroService(financeiroRepo)
	despesaService := services.NewDespesaService(despesaRepo)

	logger.OK("Serviços carregados")

	// 4. Controllers
	logger.Log("Inicializando controllers")

	usuarioController := controllers.NewUsuarioController(usuarioService, loginService)
	vendaController := controllers.NewVendasController(vendaService)
	authController := controllers.NewAuthController(loginService, auditoriaRepo)
	ordensServicoController := controllers.NewOrdemServicoController(ordensServicoService)
	clienteController := controllers.NewClientController(&clienteRepo)
	produtoController := controllers.NewProductController(&produtoRepo)
	auditoriaController := controllers.NewAuditoriaController(auditoriaService)
	trocaController := controllers.NewTrocaController(trocaService)

	financeiroController := controllers.NewFinanceiroController(financeiroService)
	despesaController := controllers.NewDespesaController(despesaService)

	logger.OK("Controllers carregados")

	// 5. Servidor HTTP e Middlewares Globais
	logger.Log("Configurando servidor HTTP")

	engine := gin.Default()

	logger.Log("Registrando middlewares")

	engine.Use(middlewares.CORSMiddleware())
	engine.Use(middlewares.JsonOnlyMiddleware())
	// engine.Use(middlewares.JsonNotEmptyMiddleware()) -> kiam: ta me dando muito trabalho!
	engine.Use(middlewares.LoggerBodyMiddleware())

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
	routes.OrdensServicoRoutes(private, ordensServicoController)
	routes.ClienteRoutes(private, clienteController)
	routes.ProdutoRoutes(private, produtoController)
	routes.VendaRoutes(private, vendaController)
	routes.AuditoriaRoute(private, auditoriaController)
	routes.RegistrarRotasTroca(private, trocaController)
	routes.FinanceiroRoutes(private, *financeiroController)
	routes.RegistrarDespesas(private, *despesaController)

	logger.OK("Rotas registradas com sucesso")

	logger.OK("Swagger disponível")
	logger.Log("Servidor iniciando na porta " + cfg.Porta)

	if err := engine.Run("0.0.0.0:" + cfg.Porta); err != nil {
		logger.Err("Falha ao iniciar servidor: " + err.Error())
		panic(err)
	}
}

package main

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/services"
	"example/web-service-gin/internal/workers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	omieService := services.NewOmieService(
		os.Getenv("OMIE_APP_KEY"),
		os.Getenv("OMIE_APP_SECRET"),
		os.Getenv("OMIE_BASE_URL"),
	)
	db := database.ConnectToDB()

	if db == nil {
		log.Fatal("Failed to connect to database")
	} else {
		log.Default().Println("Connected to database successfully")
	}
	servicoHandler := handlers.NewServicoHandler(omieService)
	clienteHandler := handlers.NewClienteHandler(omieService)
	ordemServicoHandler := handlers.NewOrdemServicoHandler(omieService)
	contaReceberHandler := handlers.NewContaReceberHandler(omieService)
	contaCorrenteHandler := handlers.NewContaCorrenteHandler(omieService)

	workerPool := workers.NewWebhookWorkerPool(omieService, 5, 100)
	defer workerPool.Shutdown()
	webhookHandler := handlers.NewWebhookHandler(workerPool)
	faturamentoCompletoHandler := handlers.NewFaturamentoCompletoHandler(omieService)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/cadastrarServico", servicoHandler.CadastrarServico)
	router.GET("/listarServicos", servicoHandler.ListarServicos)

	router.POST("/cadastrarCliente", clienteHandler.CadastrarCliente)
	router.GET("/listarClientes", clienteHandler.ListarClientes)
	router.POST("/importarCliente", clienteHandler.ImportarEmpresa)

	router.GET("/listarOrdemServico", ordemServicoHandler.ListarOrdemServicos)
	router.POST("/criarOrdemServico", ordemServicoHandler.CriarOrdemServico)
	router.POST("/faturarOrdemServico", ordemServicoHandler.FaturarOrdemServico)
	router.GET("/consultaOsFase", ordemServicoHandler.ConsultarOsFase)
	router.POST("/verificarOsFaturada", ordemServicoHandler.VerificaOsFaturada)

	router.GET("/listarContasReceber", contaReceberHandler.ListarContasReceber)
	router.POST("/consultarContaReceber", contaReceberHandler.ConsultarConta)
	router.POST("/gerarBoletoConta", contaReceberHandler.GerarBoletoConta)
	router.POST("/consultarBoletoGerado", contaReceberHandler.ConsultarBoletoGerado)
	router.POST("/webhook", webhookHandler.ReceberWebhook)
	router.POST("/criarFaturamentoCompleto", faturamentoCompletoHandler.CriarFaturamentoCompleto)
	router.POST("/consultarCliente", clienteHandler.ConsultarCliente)
	router.GET("/listarContasCorrente", contaCorrenteHandler.ListarContasCorrente)

	router.GET("/sincronizarClientes", clienteHandler.SincronizarClientes)

	router.POST("/buscarNfse", contaReceberHandler.ConsultarNFSEGerada)

	router.GET("/extratoInterCompleto", contaCorrenteHandler.ExtratoCompleto)

	fmt.Println("Rodando na porta " + os.Getenv("PORT"))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Fatal(err)
	}
}

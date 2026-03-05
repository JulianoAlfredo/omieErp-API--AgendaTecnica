package main

import (
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/services"
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
		// "7273681392978", PRODUCAO
		// "1effda944135f315ade14bdd2e7a896c", PRODUCAO
		os.Getenv("OMIE_APP_KEY"),
		os.Getenv("OMIE_APP_SECRET"),
		os.Getenv("OMIE_BASE_URL"),
	)

	servicoHandler := handlers.NewServicoHandler(omieService)
	clienteHandler := handlers.NewClienteHandler(omieService)
	ordemServicoHandler := handlers.NewOrdemServicoHandler(omieService)
	contaReceberHandler := handlers.NewContaReceberHandler(omieService)
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/cadastrarServico", servicoHandler.CadastrarServico)
	router.GET("/listarServicos", servicoHandler.ListarServicos)

	router.POST("/cadastrarCliente", clienteHandler.CadastrarCliente)
	router.GET("/listarClientes", clienteHandler.ListarClientes)

	router.GET("/listarOrdemServico", ordemServicoHandler.ListarOrdemServicos)
	router.POST("/criarOrdemServico", ordemServicoHandler.CriarOrdemServico)
	router.POST("/faturarOrdemServico", ordemServicoHandler.FaturarOrdemServico)

	router.GET("/listarContasReceber", contaReceberHandler.ListarContasReceber)
	router.POST("/consultarContaReceber", contaReceberHandler.ConsultarConta)
	router.POST("/gerarBoletoConta", contaReceberHandler.GerarBoletoConta)

	fmt.Println("Rodando na porta 8080")
	fmt.Println(os.Getenv("PORT"), os.Getenv("OMIE_APP_KEY"))
	port := os.Getenv("PORT")

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatal(err)
	}
}

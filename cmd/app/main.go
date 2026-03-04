package main

import (
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/services"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	omieService := services.NewOmieService(
		// "7273681392978", PRODUCAO
		// "1effda944135f315ade14bdd2e7a896c", PRODUCAO
		"7299234367425",
		"6de960145c93b18dc08dff314b23dfd9",
		"https://app.omie.com.br",
	)

	servicoHandler := handlers.NewServicoHandler(omieService)
	clienteHandler := handlers.NewClienteHandler(omieService)
	ordemServicoHandler := handlers.NewOrdemServicoHandler(omieService)
	router := gin.Default()

	// router.GET("/albums", func(c *gin.Context) {
	// 	c.JSON(http.StatusOK, gin.H{"message": "Albums endpoint"})
	// })

	router.POST("/cadastrarServico", servicoHandler.CadastrarServico)
	router.GET("/listarServicos", servicoHandler.ListarServicos)

	router.POST("/cadastrarCliente", clienteHandler.CadastrarCliente)
	router.GET("/listarClientes", clienteHandler.ListarClientes)

	router.POST("/criarOrdemServico", ordemServicoHandler.CriarOrdemServico)
	router.POST("/faturarOrdemServico", ordemServicoHandler.FaturarOrdemServico)


	fmt.Println("Rodando na porta 8080")
	router.Run("localhost:8080")
}

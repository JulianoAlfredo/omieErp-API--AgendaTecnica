package main

import (
	"example/web-service-gin/internal/handlers"
	"example/web-service-gin/internal/services"
	"fmt"
  "log"
  "net"
	"github.com/gin-gonic/gin"
)
func getLocalIPs() []string {
    var ips []string

    ifaces, err := net.Interfaces()
    if err != nil {
        return ips
    }

    for _, iface := range ifaces {
        // ignora interface down e loopback
        if (iface.Flags&net.FlagUp) == 0 || (iface.Flags&net.FlagLoopback) != 0 {
            continue
        }

        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addrs {
            var ip net.IP
            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }
            if ip == nil || ip.IsLoopback() {
                continue
            }

            ip = ip.To4()
            if ip == nil {
                continue // só IPv4
            }

            ips = append(ips, ip.String())
        }
    }

    return ips
}
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
	// ...existing code...
    port := "8080"

    ips := getLocalIPs()
    for _, ip := range ips {
        fmt.Printf("API disponível em: http://%s:%s\n", ip, port)
    }

    // ex.: gin
    if err := router.Run("0.0.0.0:" + port); err != nil {
        log.Fatal(err)
    }
}

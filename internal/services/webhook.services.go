package services

import (
	"example/web-service-gin/internal/models"
	"fmt"
	"net/http"
	"time"
)

func (s *OmieService) ProcessarWebhookOsFaturada(data models.WebhookOsFaturadaResponse) (int, error) {
	fmt.Printf("Processando webhook de OS faturada: Número OS: %s, Código Integração: %s, ID OS: %d\n", data.NumeroOS, data.CodigoIntegra, data.IdOs)
	time.Sleep(time.Second * 1)

	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookContaReceber(data models.WebhookContaReceberResponseInclude) (int, error) {
	fmt.Printf("Processando webhook de conta a receber incluída: Código Cliente: %f, Código Conta: %f, Número Documento: %s, Número Documento Fiscal: %s\n",
		data.CodigoCliente, data.CodigoConta, data.NumeroDocumento, data.NumeroDocumentoFiscal)

	time.Sleep(time.Second * 1)
	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookOsIncluida(data models.WebhookOsIncluidaResponse) (int, error) {
	fmt.Printf("Processando webhook de OS incluída: Código Integração: %s, ID OS: %d, ID Cliente: %d\n", data.CodigoIntegra, data.IdOs, data.IdCliente)
	time.Sleep(time.Second * 1)
	return http.StatusOK, nil
}

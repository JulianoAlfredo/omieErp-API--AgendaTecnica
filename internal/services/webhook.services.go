package services

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repositories"
	"fmt"
	"net/http"
	"time"
)

func (s *OmieService) ProcessarWebhookOsFaturada(data models.WebhookOsFaturadaResponse) (int, error) {
	fmt.Printf("Processando webhook de OS faturada: Número OS: %s, Código Integração: %s, ID OS: %d\n", data.NumeroOS, data.CodigoIntegra, data.IdOs)
	dbUpdt, err := repositories.WebhookUpdateOsFaturada(database.ConnectToDB(), fmt.Sprintf("%d", data.IdOs), fmt.Sprintf("%d", data.CodigoIntegra))
	if err != nil {
		fmt.Printf("Erro ao atualizar OS faturada: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbUpdt.RowsAffected()
		fmt.Printf("\033[32mOS faturada atualizada com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
	}
	time.Sleep(time.Second * 1)

	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookContaReceber(data models.WebhookContaReceberResponseInclude) (int, error) {
	fmt.Printf("Processando webhook de conta a receber incluída: Código Cliente: %s, Código Conta: %f, Número Documento: %s, Número Documento Fiscal: %s\n",
		data.CodigoCliente, data.CodigoConta, data.NumeroDocumento, data.NumeroDocumentoFiscal)
	time.Sleep(time.Second * 1)
	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookOsIncluida(data models.WebhookOsIncluidaResponse) (int, error) {
	db := database.ConnectToDB()

	fmt.Printf("Processando webhook de OS incluída: Código Integração: %s, ID OS: %d, ID Cliente: %d\n", data.CodigoIntegra, data.IdOs, data.IdCliente)
	dbUpdt, err := repositories.WebhookUpdateOsIncluida(db, fmt.Sprintf("%d", data.IdOs), data.CodigoIntegra)
	if err != nil {
		fmt.Printf("Erro ao atualizar OS incluída: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbUpdt.RowsAffected()
		fmt.Printf("\033[32mOS incluída atualizada com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
	}
	return http.StatusOK, nil
}

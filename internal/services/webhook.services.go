package services

import (
	"example/web-service-gin/internal/database"
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/repositories"
	"fmt"
	"net/http"
)

func (s *OmieService) ProcessarWebhookBoletoGerado(data models.WebhookBoletoGeradoResponse) (int, error) {
	fmt.Printf("Processando webhook de boleto gerado: Código Cliente: %d, Código Conta: %d, Número Pedido: %s, Boleto Gerado: %s, Código de Barras: %s, Boleto Número: %s\n",
		data.CodigoCliente, data.CodigoConta, data.NumeroPedido, data.BoletoGerado, data.CodigoBarras, data.BoletoNumero)
	dbInsertBoletoGerado, err := repositories.WebhookInsertBoletoGerado(database.ConnectToDB(), data.CodigoCliente, data.CodigoConta, data.NumeroPedido, data.BoletoGerado, data.CodigoBarras, data.BoletoNumero)
	if err != nil {
		fmt.Printf("Erro ao inserir boleto gerado: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbInsertBoletoGerado.RowsAffected()
		fmt.Printf("\033[32mBoleto gerado inserido com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)

		s.ConsultarBoletoGerado(models.ConsultaBoletoGerado{NCodTitulo: data.CodigoConta})

	}
	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookOsFaturada(data models.WebhookOsFaturadaResponse) (int, error) {
	fmt.Printf("Processando webhook de OS faturada: Número OS: %s, Código Integração: %s, ID OS: %d\n", data.NumeroOS, data.CodigoIntegra, data.IdOs)
	dbUpdt, err := repositories.WebhookUpdateOsFaturada(database.ConnectToDB(), fmt.Sprintf("%d", data.IdOs), fmt.Sprintf("%d", data.CodigoIntegra))
	if err != nil {
		fmt.Printf("Erro ao atualizar OS faturada: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbUpdt.RowsAffected()
		fmt.Printf("\033[32mOS faturada atualizada com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
	}

	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookContaReceber(data models.WebhookContaReceberResponseInclude) (int, error) {
	fmt.Printf("Processando webhook de conta a receber incluída: Código Cliente: %s, Código Conta: %f, Número Documento: %s, Número Documento Fiscal: %s, Número Pedido: %s\n",
		data.CodigoCliente, data.CodigoConta, data.NumeroDocumento, data.NumeroDocumentoFiscal, data.NumeroPedido)
	dbInsertContaReceber, err := repositories.WebhookInsertContaReceber(database.ConnectToDB(), data.CodigoCliente, data.CodigoConta, data.NumeroDocumento, data.NumeroDocumentoFiscal, data.NumeroPedido)
	if err != nil {
		fmt.Printf("Erro ao inserir conta a receber: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbInsertContaReceber.RowsAffected()
		fmt.Printf("\033[32mConta a receber inserida com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
		dbInsertDadosNFSE, err := s.UpsertNFSEGerada(models.ConsultaNFSEGerada{NNumeroNFSe: data.NumeroDocumentoFiscal})
		if err != nil {
			fmt.Printf("Erro ao obter dados da NFSe gerada: %s\n", err.Error())
		} else {
			fmt.Printf("\033[32mNFSe gerada inserida/atualizada com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
		}
		fmt.Printf("Dados da NFSe gerada: %v\n", dbInsertDadosNFSE)
	}

	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookOsIncluida(data models.WebhookOsIncluidaResponse) (int, error) {
	db := database.ConnectToDB()

	fmt.Printf("Processando webhook de OS incluída: Código Integração: %s, ID OS: %d, ID Cliente: %d, Número OS: %s\n", data.CodigoIntegra, data.IdOs, data.IdCliente, data.NumeroOs)
	dbUpdt, err := repositories.WebhookUpdateOsIncluida(db, fmt.Sprintf("%d", data.IdOs), data.CodigoIntegra, data.NumeroOs)
	if err != nil {
		fmt.Printf("Erro ao atualizar OS incluída: %s\n", err.Error())
	} else {
		rowsAffected, _ := dbUpdt.RowsAffected()
		fmt.Printf("\033[32mOS incluída atualizada com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
	}
	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookNfseAutorizada(data models.WebhookNfseAutorizadaResponse) (int, error) {
	fmt.Printf("Processando webhook de NFSe autorizada: Número OS: %s, Número RPS: %s, Código OS: %s, Código NF: %s, Data Emissão: %s\n", data.NumeroOs, data.NumeroRps, data.CodigoOs, data.CodigoNf, data.DataEmissao)
	return http.StatusOK, nil
}

func (s *OmieService) ProcessarWebhookBaixaRealizada(data models.WebhookBaixaRealizadaResponse) (int, error) {
	fmt.Printf("Processando webhook de baixa realizada: Código Lançamento Omie: %d, Código Cliente: %d\n", data.CodigoLancamentoOmie, data.CodigoCliente)
	result, err := repositories.WebhookUpdateConferido(database.ConnectToDB(), data.CodigoLancamentoOmie)
	if err != nil {
		fmt.Printf("Erro ao atualizar conferido: %s\n", err.Error())
	} else {
		rowsAffected, _ := result.RowsAffected()
		fmt.Printf("\033[32mConferido atualizado com sucesso. Linhas afetadas: %d\033[0m\n", rowsAffected)
	}
	return http.StatusOK, nil
}

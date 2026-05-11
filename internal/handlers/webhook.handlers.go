package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"example/web-service-gin/internal/workers"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	workerPool *workers.WebhookWorkerPool
}

func NewWebhookHandler(workerPool *workers.WebhookWorkerPool) *WebhookHandler {
	return &WebhookHandler{workerPool: workerPool}
}

func (h *WebhookHandler) ReceberWebhook(c *gin.Context) {
	var body map[string]interface{}
	var responseOsFaturada models.WebhookOsFaturadaResponse
	var responseContaReceber models.WebhookContaReceberResponseInclude
	var responseOsIncluida models.WebhookOsIncluidaResponse
	var responseBoletoGerado models.WebhookBoletoGeradoResponse
	var responseNfseAutorizada models.WebhookNfseAutorizadaResponse
	var responseBaixaRealizada models.WebhookBaixaRealizadaResponse

	fmt.Println("SIMULA salvamento banco: ", map[string]interface{}{
		"tipo":      body["topic"],
		"timestamp": time.Now(),
		"body":      body,
	})

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}
	fmt.Println("Webhook recebido", body["topic"])
	switch body["topic"] {

	case "Financas.ContaReceber.Incluido":
		fmt.Println("Conta a receber incluída")
		codigoCliente := body["event"].(map[string]interface{})["codigo_cliente_fornecedor"]
		codigoConta := body["event"].(map[string]interface{})["codigo_lancamento_omie"]
		numeroDocumento := body["event"].(map[string]interface{})["numero_documento"]
		numeroDocumentoFiscal := body["event"].(map[string]interface{})["numero_documento_fiscal"]
		numeroPedido := body["event"].(map[string]interface{})["numero_pedido"]

		responseContaReceber.CodigoCliente = int64(codigoCliente.(float64))
		responseContaReceber.CodigoConta = int64(codigoConta.(float64))
		responseContaReceber.NumeroDocumento = numeroDocumento.(string)
		responseContaReceber.NumeroDocumentoFiscal = numeroDocumentoFiscal.(string)
		responseContaReceber.NumeroPedido = numeroPedido.(string)
		services.GetOrquestrador().NotificarContaReceber(responseContaReceber)
		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:         workers.JobContaReceber,
			ContaReceber: &responseContaReceber,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})
	case "OrdemServico.Faturada":
		fmt.Println("Ordem de serviço faturada")
		codigoIntegra := body["event"].(map[string]interface{})["codigoIntegracao"]
		numeroOs := body["event"].(map[string]interface{})["numeroOrdemServico"]
		idOs := body["event"].(map[string]interface{})["idOrdemServico"]

		responseOsFaturada.CodigoIntegra = fmt.Sprintf("%v", codigoIntegra)
		responseOsFaturada.NumeroOS = fmt.Sprintf("%v", numeroOs)
		responseOsFaturada.IdOs = int64(idOs.(float64))

		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:       workers.JobOsFaturada,
			OsFaturada: &responseOsFaturada,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})
	case "OrdemServico.Incluida":
		fmt.Println("Ordem de serviço incluída")
		codigoIntegra := body["event"].(map[string]interface{})["codigoIntegracao"]
		idOs := body["event"].(map[string]interface{})["idOrdemServico"]
		idCliente := body["event"].(map[string]interface{})["idCliente"]
		numeroOs := body["event"].(map[string]interface{})["numeroOrdemServico"]

		responseOsIncluida.CodigoIntegra = fmt.Sprintf("%v", codigoIntegra)
		responseOsIncluida.IdOs = int64(idOs.(float64))
		responseOsIncluida.IdCliente = int64(idCliente.(float64))
		responseOsIncluida.NumeroOs = fmt.Sprintf("%v", numeroOs)
		services.GetOrquestrador().NotificarOsIncluida(responseOsIncluida)
		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:       workers.JobOsIncluida,
			OsIncluida: &responseOsIncluida,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})

	case "OrdemServico.Alterada":
		fmt.Println("Ordem de serviço alterada")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de ordem de serviço alterada recebido"})
	case "NFSe.NotaAutorizada":
		fmt.Println("Nota fiscal de serviço autorizada")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de nota fiscal de serviço autorizada recebido"})

		xmlNfe := body["event"].(map[string]interface{})["nfse_xml"]
		numeroOs := body["event"].(map[string]interface{})["numero_os"]
		numeroRps := body["event"].(map[string]interface{})["numero_rps"]
		codigoOs := body["event"].(map[string]interface{})["codigo_os"]
		codigoNf := body["event"].(map[string]interface{})["id_nf"]
		dataEmissao := body["event"].(map[string]interface{})["data_emis"]

		responseNfseAutorizada.NumeroOs = fmt.Sprintf("%v", numeroOs)
		responseNfseAutorizada.NumeroRps = fmt.Sprintf("%v", numeroRps)
		responseNfseAutorizada.NFseXML = fmt.Sprintf("%v", xmlNfe)
		responseNfseAutorizada.CodigoOs = fmt.Sprintf("%v", codigoOs)
		responseNfseAutorizada.CodigoNf = fmt.Sprintf("%v", codigoNf)
		responseNfseAutorizada.DataEmissao = fmt.Sprintf("%v", dataEmissao)
		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:           workers.JobNfseAutorizada,
			NfseAutorizada: &responseNfseAutorizada,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})

	case "OrdemServico.EtapaAlterada":
		fmt.Println("Ordem de serviço etapa alterada")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de ordem de serviço etapa alterada recebido"})
	case "Financas.ContaReceber.Excluido":
		fmt.Println("Conta a receber excluída")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de conta a receber excluída recebido"})
	case "NFSe.NotaCancelada":
		fmt.Println("Nota fiscal de serviço cancelada")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de nota fiscal de serviço cancelada recebido"})
	case "NFSe.NotaSubstituida":
		fmt.Println("Nota fiscal de serviço substituída")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de nota fiscal de serviço substituída recebido"})
	case "OrdemServico.Cancelada":
		fmt.Println("Ordem de serviço cancelada")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de ordem de serviço cancelada recebido"})
	case "OrdemServico.Excluida":
		fmt.Println("Ordem de serviço excluída")
		c.JSON(http.StatusOK, gin.H{"message": "Webhook de ordem de serviço excluída recebido"})
	case "Financas.ContaReceber.BoletoGerado":
		fmt.Println("Boleto gerado")
		codigoConta := body["event"].(map[string]interface{})["codigo_lancamento_omie"]
		idCliente := body["event"].(map[string]interface{})["codigo_cliente_fornecedor"]
		codigoBarras := body["event"].(map[string]interface{})["codigo_barras_ficha_compensacao"]
		BoletoGerado := body["event"].(map[string]interface{})["boleto_gerado"]
		numeroOs := body["event"].(map[string]interface{})["numero_pedido"]
		BoletoNumero := body["event"].(map[string]interface{})["boleto_numero"]

		responseBoletoGerado.CodigoConta = int64(codigoConta.(float64))
		responseBoletoGerado.NumeroPedido = fmt.Sprintf("%v", numeroOs)
		responseBoletoGerado.CodigoCliente = int64(idCliente.(float64))
		responseBoletoGerado.BoletoGerado = fmt.Sprintf("%v", BoletoGerado)
		responseBoletoGerado.CodigoBarras = fmt.Sprintf("%v", codigoBarras)
		responseBoletoGerado.BoletoNumero = fmt.Sprintf("%v", BoletoNumero)
		services.GetOrquestrador().NotificarBoletoGerado(responseBoletoGerado)
		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:         workers.JobBoletoGerado,
			BoletoGerado: &responseBoletoGerado,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})
	case "Financas.ContaReceber.BaixaRealizada":
		fmt.Println("Baixa realizada em conta a receber")
		event := body["event"].([]interface{})
		if len(event) > 0 {
			eventData := event[0].(map[string]interface{})
			contaReceber := eventData["conta_a_receber"].([]interface{})
			if len(contaReceber) > 0 {
				contaData := contaReceber[0].(map[string]interface{})
				responseBaixaRealizada.CodigoLancamentoOmie = int64(contaData["codigo_lancamento_omie"].(float64))
				responseBaixaRealizada.CodigoCliente = int64(eventData["codigo_cliente_fornecedor"].(float64))
				responseBaixaRealizada.Data = fmt.Sprintf("%v", eventData["data"])
				responseBaixaRealizada.DataCred = fmt.Sprintf("%v", eventData["data_cred"])
				responseBaixaRealizada.Observacao = fmt.Sprintf("%v", eventData["observacao"])
				responseBaixaRealizada.Valor = eventData["valor"].(float64)
				err := h.workerPool.Enqueue(workers.WebhookJob{
					Tipo:           workers.JobBaixaRealizada,
					BaixaRealizada: &responseBaixaRealizada,
				})
				if err != nil {
					c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
					return
				}
			}
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Webhook recebido com sucesso"})
}

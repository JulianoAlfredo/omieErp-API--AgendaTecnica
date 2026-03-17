package handlers

import (
	"example/web-service-gin/internal/models"
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

	fmt.Println("SIMULA salvamento banco: ", map[string]interface{}{
		"tipo":      body["topic"],
		"timestamp": time.Now(),
		"body":      body,
	})
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	switch body["topic"] {
	case "Financas.ContaReceber.Incluido":
		fmt.Println("Conta a receber incluída")
		codigoCliente := body["event"].(map[string]interface{})["codigo_cliente_fornecedor"]
		codigoConta := body["event"].(map[string]interface{})["codigo_lancamento_omie"]
		numeroDocumento := body["event"].(map[string]interface{})["numero_documento"]
		numeroDocumentoFiscal := body["event"].(map[string]interface{})["numero_documento_fiscal"]

		responseContaReceber.CodigoCliente = int64(codigoCliente.(float64))
		responseContaReceber.CodigoConta = int64(codigoConta.(float64))
		responseContaReceber.NumeroDocumento = numeroDocumento.(string)
		responseContaReceber.NumeroDocumentoFiscal = numeroDocumentoFiscal.(string)
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

		responseOsIncluida.CodigoIntegra = fmt.Sprintf("%v", codigoIntegra)
		responseOsIncluida.IdOs = int64(idOs.(float64))
		responseOsIncluida.IdCliente = int64(idCliente.(float64))

		err := h.workerPool.Enqueue(workers.WebhookJob{
			Tipo:       workers.JobOsIncluida,
			OsIncluida: &responseOsIncluida,
		})
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"erro": "fila cheia, tente novamente"})
			return
		}
		c.JSON(http.StatusAccepted, gin.H{"status": "enfileirado"})

	}
	c.JSON(http.StatusOK, gin.H{"message": "Webhook recebido com sucesso"})
}

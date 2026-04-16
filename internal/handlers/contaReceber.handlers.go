package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ContaReceberHandler struct {
	omieService *services.OmieService
}

func NewContaReceberHandler(omieService *services.OmieService) *ContaReceberHandler {
	return &ContaReceberHandler{omieService: omieService}
}

func (h *ContaReceberHandler) ListarContasReceber(c *gin.Context) {

	resultado, err := h.omieService.ListarContasReceber()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

func (h *ContaReceberHandler) ConsultarConta(c *gin.Context) {
	var req models.ContaReceberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.ConsultarContaReceber(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

func (h *ContaReceberHandler) GerarBoletoConta(c *gin.Context) {
	var req models.GerarBoletoConta

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "body inválido",
			"detalhe": err.Error(),
		})
		return
	}

	resultado, err := h.omieService.GerarBoletoConta(req)
	fmt.Println(resultado)
	if resultado["cNumBoleto"] != "" && resultado["cNumBoleto"] != nil {
		c.JSON(http.StatusOK, resultado)
		fmt.Println("Boleto gerado com sucesso! Número do boleto:", resultado["cNumBoleto"])
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao gerar boleto", "cod_stautus": resultado["cod_status"], "msg_status": resultado["cDesStatus"]})

	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (h *ContaReceberHandler) ConsultarBoletoGerado(c *gin.Context) {
	var req models.ConsultaBoletoGerado
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "body inválido",
			"detalhe": err.Error(),
		})
		return
	}
	resultado, err := h.omieService.ConsultarBoletoGerado(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

func (h *ContaReceberHandler) ConsultarNFSEGerada(c *gin.Context) {
	var req models.ConsultaNFSEGerada
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "body inválido",
			"detalhe": err.Error(),
		})
		return
	}
	resultado, err := h.omieService.UpsertNFSEGerada(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

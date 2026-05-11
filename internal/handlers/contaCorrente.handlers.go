package handlers

import (
	"example/web-service-gin/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ContaCorrenteHandler struct {
	omieService *services.OmieService
}

func NewContaCorrenteHandler(omieService *services.OmieService) *ContaCorrenteHandler {
	return &ContaCorrenteHandler{omieService: omieService}
}

func (h *ContaCorrenteHandler) ListarContasCorrente(c *gin.Context) {
	resultado, err := h.omieService.ListarContasCorrente()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

func (h *ContaCorrenteHandler) ExtratoCompleto(c *gin.Context) {
	resultado, err := h.omieService.ExtratoCompleto()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

type sincronizarBaixasRequest struct {
	NCodCC          int64  `json:"nCodCC"`
	DPeriodoInicial string `json:"dPeriodoInicial"`
	DPeriodoFinal   string `json:"dPeriodoFinal"`
}

func (h *ContaCorrenteHandler) SincronizarBaixas(c *gin.Context) {
	var req sincronizarBaixasRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "body invalido: " + err.Error()})
		return
	}
	if req.NCodCC == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erro": "nCodCC e obrigatorio"})
		return
	}

	now := time.Now()
	dataInicial := req.DPeriodoInicial
	if dataInicial == "" {
		dataInicial = "01/" + now.Format("01/2006")
	}
	dataFinal := req.DPeriodoFinal
	if dataFinal == "" {
		dataFinal = now.Format("02/01/2006")
	}

	resultado, err := h.omieService.SincronizarBaixasOmie(req.NCodCC, dataInicial, dataFinal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erro": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

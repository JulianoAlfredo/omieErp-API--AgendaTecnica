package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type FaturamentoCompletoHandler struct {
	omieService *services.OmieService
}

func NewFaturamentoCompletoHandler(omieService *services.OmieService) *FaturamentoCompletoHandler {
	return &FaturamentoCompletoHandler{omieService: omieService}
}

func (h *FaturamentoCompletoHandler) CriarFaturamentoCompleto(c *gin.Context) {
	var req models.OrdemServicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.CriarFaturamentoCompleto(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultado)
}

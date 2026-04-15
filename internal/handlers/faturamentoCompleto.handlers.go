package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"io"
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
	
	progressCh := make(chan models.FaturamentoProgresso, 10)
	go h.omieService.CriarFaturamentoCompletoStream(req, progressCh)

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	c.Stream(func(w io.Writer) bool {
		if evento, ok := <-progressCh; ok {
			c.SSEvent("progresso", evento)
			return true
		}
		return false
	})
}

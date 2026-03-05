package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.GerarBoletoConta(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resultado)
}

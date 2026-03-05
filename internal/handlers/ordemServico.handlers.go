package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrdemServicoHandler struct {
	omieService *services.OmieService
}

func NewOrdemServicoHandler(omieService *services.OmieService) *OrdemServicoHandler {
	return &OrdemServicoHandler{omieService: omieService}
}

func (h *OrdemServicoHandler) CriarOrdemServico(c *gin.Context) {
	var req models.OrdemServicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.CriarOrdemServico(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resultado})
}
func (h *OrdemServicoHandler) ListarOrdemServicos(c *gin.Context) {

	resultado, err := h.omieService.ListarOrdemServico()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resultado})
}
func (h *OrdemServicoHandler) FaturarOrdemServico(c *gin.Context) {
	var req models.FaturaOrdemServicoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.FaturarOrdemServico(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resultado})
}

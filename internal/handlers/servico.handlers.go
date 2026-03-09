package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServicoHandler struct {
	omieService *services.OmieService
}

func NewServicoHandler(omieService *services.OmieService) *ServicoHandler {
	return &ServicoHandler{omieService: omieService}
}

func (h *ServicoHandler) CadastrarServico(c *gin.Context) {
	var req models.ServicoRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.CadastrarServico(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resultado})
}
func (h *ServicoHandler) ListarServicos(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Listar serviços endpoint"})

	resultado, err := h.omieService.ListarServicos()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": resultado})

}

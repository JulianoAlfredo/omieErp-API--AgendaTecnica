package handlers

import (
	"example/web-service-gin/internal/models"
	"example/web-service-gin/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ClienteHandler struct {
	omieService *services.OmieService
}

func NewClienteHandler(omieService *services.OmieService) *ClienteHandler {
	return &ClienteHandler{omieService: omieService}
}

func (h *ClienteHandler) CadastrarCliente(c *gin.Context) {
	var req models.ClienteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resultado, err := h.omieService.CriarCliente(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": resultado})
}

func (h *ClienteHandler) ListarClientes(c *gin.Context) {

	resultado, err := h.omieService.ListarClientes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resultado)

}

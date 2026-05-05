package handlers

import (
	"example/web-service-gin/internal/services"
	"net/http"

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

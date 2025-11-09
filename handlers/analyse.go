package handlers

import (
	"net/http"
	"web-analyser/models"
	"web-analyser/services"

	"github.com/gin-gonic/gin"
)

type AnalyseHandler struct {
	analyserService *services.AnalyserService
}

func NewAnalyseHandler(analyserService *services.AnalyserService) *AnalyseHandler {
	return &AnalyseHandler{analyserService: analyserService}
}

func (h *AnalyseHandler) AnalyzePage(c *gin.Context) {

	var req models.WebAnalysingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request trying to bind for web analysing request: " + err.Error(),
		})
		return
	}

}

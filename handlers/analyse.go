package handlers

import (
	"net/http"
	"web-analyser/models"
	"web-analyser/services"

	"github.com/gin-gonic/gin"
)

type AnalyseHandler struct {
	analyserService *services.AnalyserService
	workerPool      chan struct{}
}

func NewAnalyseHandler(analyserService *services.AnalyserService, maxWorkers int) *AnalyseHandler {
	return &AnalyseHandler{
		analyserService: analyserService,
		workerPool:      make(chan struct{}, maxWorkers),
	}
}

func (h *AnalyseHandler) AnalysePage(c *gin.Context) {

	var req models.WebAnalysingRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request trying to bind for web analysing request: " + err.Error(),
		})
		return
	}

	select {
	case h.workerPool <- struct{}{}:
		defer func() { <-h.workerPool }()
	default:
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Server busy, please try again later",
		})
		return
	}

	result, err := h.analyserService.AnalyserWebUrl(req.Url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.WebAnalysingResponse{
			Error: "Analysis failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)

}

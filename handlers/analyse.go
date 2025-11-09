package handlers

import (
	"fmt"
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
	fmt.Printf("Incoming analysing request through gin setup %v", c)
}

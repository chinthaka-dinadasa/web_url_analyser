package main

import (
	"net/http"
	"time"
	"web-analyser/handlers"
	"web-analyser/services"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	analyzer := services.NewHtmlAnalyser()
	analyzeHandler := handlers.NewAnalyseHandler(analyzer)

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health":      "up",
			"api":         "web-analyser-api-v1.0",
			"requestedOn": time.Now().Format(time.RFC3339),
		})
	})

	router.POST("/process-web-url", analyzeHandler.AnalyzePage)

	router.Run()
}

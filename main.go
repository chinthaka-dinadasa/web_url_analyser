package main

import (
	"net/http"
	"time"
	"web-analyser/handlers"
	"web-analyser/logger"
	"web-analyser/models"
	"web-analyser/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	logger.InitLogger()

	router := gin.Default()
	config := models.LoadConfig()
	analyser := services.NewAnalyserService(config.CacheTtl)
	analyseHandler := handlers.NewAnalyseHandler(analyser, config.MaxWorkers)

	router.Use(cors.Default()) //TODO setup cors with ENV based origin.

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health":      "up",
			"api":         "web-analyser-api-v1.0",
			"requestedOn": time.Now().Format(time.RFC3339),
		})
	})

	router.POST("/process-web-url", analyseHandler.AnalysePage)

	logger.Info("Server starting on default port", config.Port)
	router.Run(":" + config.Port)
}

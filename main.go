package main

import (
	"net/http"
	"time"
	_ "web-analyser/docs"
	"web-analyser/handlers"
	"web-analyser/logger"
	"web-analyser/models"
	"web-analyser/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Web Analyser API - Golang
// @version 1.0
// @description Web analyser API developed using Golang
// @contact.name Chinthaka D
// @host localhost:8080
func main() {
	logger.InitLogger()

	router := gin.Default()
	config := models.LoadConfig()
	analyser := services.NewAnalyserService(config.CacheTtl)
	analyseHandler := handlers.NewAnalyseHandler(analyser, config.MaxWorkers)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Use(cors.Default())

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

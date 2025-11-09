package main

import (
	"net/http"
	"time"
	"web-analyser/models"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"health":      "up",
			"api":         "web-analyser-api-v1.0",
			"requestedOn": time.Now().Format(time.RFC3339),
		})
	})

	router.POST("/process-web-url", processWebUrl)

	router.Run()
}

func processWebUrl(c *gin.Context) {
	var webAnylysingRequest models.WebAnalysingRequest
	// Bind JSON to struct with validation
	if err := c.ShouldBindJSON(&webAnylysingRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// Process user data
	c.JSON(http.StatusCreated, gin.H{
		"status": "processing",
		"weburl": webAnylysingRequest,
	})
}

package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type WebAnalysingRequest struct {
	Url string `json:"url" binding:"required"`
}

type WebAnalysingResponse struct {
	HTMLVersion string `json:"htmlVersion"`
	PageTitle   string `json:"pageTitle"`
}

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
	var webAnylysingRequest WebsiteAnalysingRequest
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

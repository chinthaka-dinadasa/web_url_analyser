package services

import (
	"fmt"
	"net/http"
	"web-analyser/models"
)

type AnalyserService struct {
	client *http.Client
}

func NewAnalyserService() *AnalyserService {
	return &AnalyserService{
		client: &http.Client{},
	}
}

func (a *AnalyserService) AnalyserWebUrl(targetURL string) (*models.WebAnalysingResponse, error) {
	resp, err := a.client.Get(targetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	result := &models.WebAnalysingResponse{
		HTMLVersion: "5.0.0",
		PageTitle:   "Hello",
	}

	return result, nil
}

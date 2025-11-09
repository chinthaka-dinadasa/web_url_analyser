package services

import (
	"fmt"
	"net/http"
	"web-analyser/models"

	"github.com/PuerkitoBio/goquery"
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

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	result := &models.WebAnalysingResponse{
		HTMLVersion: a.captureHTMLVersion(doc),
		PageTitle:   a.capturePageTitle(doc),
	}

	return result, nil
}

func (a *AnalyserService) captureHTMLVersion(doc *goquery.Document) string {
	//panic("unimplemented")
	return "HTML5"
}

func (a *AnalyserService) capturePageTitle(doc *goquery.Document) string {
	//panic("unimplemented")
	return "Javatodev.com"
}

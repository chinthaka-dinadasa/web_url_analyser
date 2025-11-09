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

	numberOfInternalLinks, unaccessibleInLinks := a.captureInternalLinks(doc)
	numberOfExternalLinks, unaccessibleExLinks := a.captureExternalLinks(doc)

	result := &models.WebAnalysingResponse{
		HTMLVersion:       a.captureHTMLVersion(doc),
		PageTitle:         a.capturePageTitle(doc),
		Heading:           a.captureHeadingDetails(doc),
		InternalLinks:     numberOfInternalLinks,
		ExternalLinks:     numberOfExternalLinks,
		UnaccessibleLinks: unaccessibleExLinks + unaccessibleInLinks,
	}

	return result, nil
}

func (a *AnalyserService) captureExternalLinks(doc *goquery.Document) (int16, int16) {
	return 8, 1
}

func (a *AnalyserService) captureInternalLinks(doc *goquery.Document) (int16, int16) {
	return 8, 0
}

func (a *AnalyserService) captureHeadingDetails(doc *goquery.Document) models.HeadingDetail {
	return models.HeadingDetail{
		H1: 2,
		H2: 3,
		H3: 4,
		H4: 0,
		H5: 2,
		H6: 0,
	}
}

func (a *AnalyserService) captureHTMLVersion(doc *goquery.Document) string {

	return "HTML5"
}

func (a *AnalyserService) capturePageTitle(doc *goquery.Document) string {
	//panic("unimplemented")
	return doc.Find("title").Text()
}

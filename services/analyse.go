package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"web-analyser/logger"
	"web-analyser/models"

	"github.com/PuerkitoBio/goquery"
)

type AnalyserService struct {
	client *http.Client
	cache  *SimpleCache
}

func NewAnalyserService() *AnalyserService {
	cache := NewSimpleCache()
	return &AnalyserService{
		client: &http.Client{},
		cache:  cache,
	}
}

func (a *AnalyserService) AnalyserWebUrl(targetURL string) (*models.WebAnalysingResponse, error) {

	logger.Info("Analysing......", "url", targetURL)

	data := a.cache.Get(targetURL)
	if data != nil {
		logger.Debug("Cache hit", "url", targetURL)
		return data, nil
	}

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
		HTMLVersion:           a.captureHTMLVersion(doc),
		PageTitle:             a.capturePageTitle(doc),
		Heading:               a.captureHeadingDetails(doc),
		LinksData:             a.captureLinksData(targetURL, doc),
		LoginFormAvailability: a.captureLoginForm(doc),
	}

	a.cache.Set(targetURL, result)

	return result, nil
}

func (a *AnalyserService) captureLoginForm(doc *goquery.Document) bool {

	if doc.Find("form input[type='password' i]").Length() > 0 {
		return true
	}

	foundLoginFormData := false

	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		if form.Find("input[type*='password' i]").Length() > 0 {
			foundLoginFormData = true
		}
	})
	// TODO: add checks for login form with Login sign in button texts if time permits
	return foundLoginFormData

}

func (a *AnalyserService) captureLinksData(baseUrl string, doc *goquery.Document) models.WebLinkDetail {

	base, err := url.Parse(baseUrl)

	var webLinkDetails models.WebLinkDetail
	if err != nil {
		logger.Error("External Link capturing failed", "err", err)
	} else {
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			href, _ := s.Attr("href")

			linkUrl, err := url.Parse(href)
			if err != nil {
				logger.Error("Error in parsing link", "link", href, "err", err)
				return
			}

			if !linkUrl.IsAbs() {
				linkUrl = base.ResolveReference(linkUrl)
			}

			if linkUrl.Host == base.Host {
				webLinkDetails.InternalLinks++
			} else if a.validLinkUrl(linkUrl) {
				webLinkDetails.ExternalLinks++
				if !a.isLinkAccessible(linkUrl.String()) {
					webLinkDetails.UnAccessibleLinks++
				}
			}

		})
	}

	return webLinkDetails
}

func (a *AnalyserService) validLinkUrl(linkUrl *url.URL) bool {
	invalidPrefixes := []string{
		"mailto:", "javascript:", "tel:", //TODO add more prefixes here.
	}
	for _, prefix := range invalidPrefixes {
		if strings.Contains(strings.ToLower(linkUrl.String()), prefix) {
			return false
		}
	}
	return true
}

func (a *AnalyserService) isLinkAccessible(link string) bool {
	_, err := a.client.Get(link)
	return err == nil
}

func (a *AnalyserService) captureHeadingDetails(doc *goquery.Document) models.HeadingDetail {
	return models.HeadingDetail{
		H1: doc.Find("h1").Length(),
		H2: doc.Find("h2").Length(),
		H3: doc.Find("h3").Length(),
		H4: doc.Find("h4").Length(),
		H5: doc.Find("h5").Length(),
		H6: doc.Find("h6").Length(),
	}
}

func (a *AnalyserService) captureHTMLVersion(doc *goquery.Document) string {
	htmlContent, err := doc.Html()
	if err != nil {
		logger.Error("Error in parsing docs HTML content for HTML version capturing", "err", err)
	}

	var decodedHtmlContent string

	if err := json.Unmarshal([]byte(`"`+htmlContent+`"`), &decodedHtmlContent); err == nil {
		htmlContent = decodedHtmlContent
	}

	html := strings.TrimSpace(strings.ToLower(htmlContent))

	switch {
	case strings.Contains(html, "html 4.01"):
		return "HTML4"
	case strings.Contains(html, "xhtml"):
		return "XHTML"
	case strings.Contains(html, "<!doctype html"):
		return "HTML5"
	case strings.Contains(html, "<html"):
		return "HTML"
	default:
		return "UNIDENTIFIED"
	}
}

func (a *AnalyserService) capturePageTitle(doc *goquery.Document) string {
	return doc.Find("title").Text()
}

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"web-analyser/logger"
	"web-analyser/models"

	"github.com/PuerkitoBio/goquery"
)

type AnalyserService struct {
	client *http.Client
	cache  *SimpleCache
}

func NewAnalyserService(ttl int) *AnalyserService {
	cache := NewSimpleCache(ttl)
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
		externalLinks := []string{}
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
				externalLinks = append(externalLinks, linkUrl.String())
			}

		})

		if len(externalLinks) > 0 {
			accessibilityResults := a.checkExternalLinkAccessibility(externalLinks)
			for _, accessible := range accessibilityResults {
				if !accessible {
					webLinkDetails.UnAccessibleLinks++
				}
			}
		}

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

	maxRetries := 3
	for attmpt := 1; attmpt <= maxRetries; attmpt++ {
		resp, err := a.client.Get(link)
		if err == nil && resp != nil {
			defer resp.Body.Close()
			if resp.StatusCode < 400 {
				fmt.Printf("Success: GET URL %v - %v attempt %d\n", link, resp.StatusCode, attmpt)
				return true
			}
		} else {
			fmt.Printf("Failed: GET URL %v - attempt %d\n", link, attmpt)
		}

		if attmpt < maxRetries {
			fmt.Printf("Retrying: GET URL %v failed (attempt %d) - Error: %v\n", link, attmpt, err)
		}
	}
	fmt.Printf("Failed: GET URL %v after %d attempts\n", link, maxRetries)
	return false
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

func (a *AnalyserService) checkExternalLinkAccessibility(links []string) map[string]bool {
	results := make(map[string]bool)
	var mutex sync.Mutex
	var wg sync.WaitGroup

	numberOfWorkrs := 5
	linkScanChannel := make(chan string, len(links))

	for range numberOfWorkrs {
		wg.Go(func() {
			for link := range linkScanChannel {
				isAccessible := a.isLinkAccessible(link)
				mutex.Lock()

				results[link] = isAccessible
				mutex.Unlock()

			}
		})
	}

	for _, link := range links {
		linkScanChannel <- link
	}

	close(linkScanChannel)
	wg.Wait()
	return results
}

package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

	webLinksData := a.captureLinksData(targetURL, doc)

	result := &models.WebAnalysingResponse{
		HTMLVersion: a.captureHTMLVersion(doc),
		PageTitle:   a.capturePageTitle(doc),
		Heading:     a.captureHeadingDetails(doc),
		LinksData:   webLinksData,
	}

	return result, nil
}

func (a *AnalyserService) captureLinksData(baseUrl string, doc *goquery.Document) models.WebLinkDetail {

	base, err := url.Parse(baseUrl)
	if err != nil {
		fmt.Printf("External link capturing failed %v", err)
	}
	fmt.Printf("Base url %v\n", base.Host)
	var webLinkDetails models.WebLinkDetail
	doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")

		linkUrl, err := url.Parse(href)
		if err != nil {
			fmt.Printf("Error in parsing link %v \n", err)
			return
		}

		if !linkUrl.IsAbs() {
			linkUrl = base.ResolveReference(linkUrl)
		}
		fmt.Printf("Href url %v\n", linkUrl)
		if linkUrl.Host == base.Host {
			webLinkDetails.UnAccessibleLinks++
		} else {
			webLinkDetails.ExternalLinks++
			if !a.isLinkAccessible(linkUrl.String()) {
				webLinkDetails.UnAccessibleLinks++
			}
		}

	})
	return webLinkDetails
}

func (a *AnalyserService) isLinkAccessible(link string) bool {
	_, err := a.client.Get(link)
	fmt.Printf("Data comming from link accessibilty test %v", err)
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
		fmt.Printf("Error in parsing docs HTML content for HTML version capturing %v \n", err)
	}

	var decodedHtmlContent string

	if err := json.Unmarshal([]byte(`"`+htmlContent+`"`), &decodedHtmlContent); err == nil {
		htmlContent = decodedHtmlContent
	}

	html := strings.TrimSpace(strings.ToLower(htmlContent))

	switch {
	case strings.Contains(html, "<!doctype html>"):
		return "HTML5"
	case strings.Contains(html, "html 4.01"):
		return "HTML4"
	case strings.Contains(html, "xhtml"):
		return "XHTML"
	case strings.Contains(html, "<html"):
		return "HTML"
	default:
		return "UNIDENTIFIED"
	}
}

func (a *AnalyserService) capturePageTitle(doc *goquery.Document) string {
	return doc.Find("title").Text()
}

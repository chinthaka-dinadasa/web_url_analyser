package services

import (
	"strings"
	"testing"
	"web-analyser/models"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestAnalyserService_PageTitleExtraction(t *testing.T) {
	analyser := NewAnalyserService()
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "basic title",
			html:     `<html><head><title>Test Page</title></head><body></body></html>`,
			expected: "Test Page",
		},
		{
			name:     "basic title",
			html:     `<title>Random</title>`,
			expected: "Random",
		},
		{
			name:     "basic title",
			html:     `<html>hello</html>`,
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.capturePageTitle(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureLoginForm(t *testing.T) {
	analyser := NewAnalyserService()
	tests := []struct {
		name     string
		html     string
		expected bool
	}{
		{
			name:     "simple password input in form",
			html:     `<form><input type="password" name="password"></form>`,
			expected: true,
		},
		{
			name:     "password input with different attributes",
			html:     `<form><input type="password" placeholder="Enter password"></form>`,
			expected: true,
		},
		{
			name:     "multiple password inputs",
			html:     `<form><input type="password" name="pwd1"><input type="password" name="pwd2"></form>`,
			expected: true,
		},
		{
			name:     "password input with other inputs",
			html:     `<form><input type="text" name="username"><input type="password" name="password"></form>`,
			expected: true,
		},
		{
			name:     "form without password input",
			html:     `<form><input type="text" name="username"><input type="email" name="email"></form>`,
			expected: false,
		},
		{
			name:     "empty form",
			html:     `<form></form>`,
			expected: false,
		},
		{
			name:     "no forms at all",
			html:     `<html><body><div>Hello World</div></body></html>`,
			expected: false,
		},
		{
			name:     "password input outside form",
			html:     `<input type="password"><form><input type="text"></form>`,
			expected: false,
		},
		{
			name: "multiple forms with one having password",
			html: `
            <form id="search">
                <input type="text" placeholder="Search">
            </form>
            <form id="login">
                <input type="text" name="username">
                <input type="password" name="password">
            </form>
            <form id="contact">
                <input type="text" name="name">
            </form>`,
			expected: true,
		},
		{
			name:     "nested password input",
			html:     `<form><div><section><input type="password" name="pwd"></section></div></form>`,
			expected: true,
		},
		{
			name:     "password input with mixed case",
			html:     `<form><input type="PASSWORD" name="password"></form>`,
			expected: true,
		},
		{
			name:     "malformed password input",
			html:     `<form><input type="password" ></form>`,
			expected: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.captureLoginForm(doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyserService_CaptureLinksData(t *testing.T) {
	analyser := NewAnalyserService()
	tests := []struct {
		name     string
		html     string
		baseUrl  string
		expected models.WebLinkDetail
	}{
		{
			name: "Internal links",
			html: `<html>
		        <body>
		            <a href="/about">About</a>
		            <a href="/contact">Contact</a>
		            <a href="#section">Anchor</a>
		        </body>
		    </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     3,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "Internal and external links",
			html: `<html>
		        <body>
		            <a href="/about">About</a>
		            <a href="/contact">Contact</a>
		            <a href="https://github.com/">Anchor</a>
		        </body>
		    </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     2,
				ExternalLinks:     1,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "Internal and external links with 404",
			html: `<html>
                <body>
                    <a href="/about">About</a>
                    <a href="/contact">Contact</a>
                    <a href="https://noname_757a971d-ac55-4651-8622-17e62b703310393.com">Anchor</a>
                </body>
            </html>`,
			baseUrl: "https://www.javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     2,
				ExternalLinks:     1,
				UnAccessibleLinks: 1,
			},
		},
		{
			name:    "empty page",
			html:    `<html><body>No links here</body></html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     0,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "links without href",
			html: `
            <html>
                <body>
                    <a name="anchor">No href</a>
                    <a>Empty anchor</a>
                    <a href="/valid">Valid link</a>
                </body>
            </html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     1,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
		{
			name: "mailto and javascript links",
			html: `
            <html>
                <body>
                    <a href="mailto:author@javatodev.com">Email</a>
                    <a href="javascript:void(0)">JS Link</a>
                    <a href="tel:+1234567890">Phone</a>
                    <a href="/normal">Normal Link</a>
                </body>
            </html>`,
			baseUrl: "https://javatodev.com",
			expected: models.WebLinkDetail{
				InternalLinks:     1,
				ExternalLinks:     0,
				UnAccessibleLinks: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err)

			result := analyser.captureLinksData(tt.baseUrl, doc)
			assert.Equal(t, tt.expected, result)
		})
	}
}

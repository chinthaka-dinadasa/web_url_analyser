package services

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestHtmlAnalyser_PageTitleExtraction(t *testing.T) {
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

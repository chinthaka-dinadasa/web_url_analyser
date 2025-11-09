package services

import "net/http"

type AnalyserService struct {
	client *http.Client
}

func NewHtmlAnalyser() *AnalyserService {
	return &AnalyserService{
		client: &http.Client{},
	}
}

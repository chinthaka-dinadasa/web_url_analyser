package handlers

import "net/http"

type HTMLAnlyser struct {
	client *http.Client
}

func NewHtmlAnalyser() *HTMLAnlyser {
	return &HTMLAnlyser{
		client: &http.Client{},
	}
}

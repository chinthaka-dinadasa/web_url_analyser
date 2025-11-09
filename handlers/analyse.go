package handlers

import "web-analyser/services"

type AnalyseHadler struct {
	analyserService *services.AnalyserService
}

func NewAnalyseHandler(analyserService *services.AnalyserService) *AnalyseHandler {
	return &AnalyseHandler{analyserService: analyserService}
}

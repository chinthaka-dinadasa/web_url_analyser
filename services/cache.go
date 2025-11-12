package services

import (
	"web-analyser/models"
)

type SimpleCache struct {
	data map[string]*models.WebAnalysingResponse
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string]*models.WebAnalysingResponse),
	}
}

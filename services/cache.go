package services

import (
	"sync"
	"web-analyser/models"
)

type SimpleCache struct {
	data  map[string]*models.WebAnalysingResponse
	mutex sync.RWMutex
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data: make(map[string]*models.WebAnalysingResponse),
	}
}

func (c *SimpleCache) Get(url string) *models.WebAnalysingResponse {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if result, exists := c.data[url]; exists {
		return result
	}
	return nil
}

func (c *SimpleCache) Set(url string, result *models.WebAnalysingResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[url] = result
}

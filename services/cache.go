package services

import (
	"sync"
	"time"
	"web-analyser/models"
)

type SimpleCache struct {
	data  map[string]*models.WebAnalysingResponse
	mutex sync.RWMutex
	times map[string]time.Time
	ttl   time.Duration
}

func NewSimpleCache() *SimpleCache {
	return &SimpleCache{
		data:  make(map[string]*models.WebAnalysingResponse),
		times: make(map[string]time.Time),
		ttl:   1 * time.Hour, //TODO move this to ENV when adding env part
	}
}

func (c *SimpleCache) Get(url string) *models.WebAnalysingResponse {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if result, exists := c.data[url]; exists {
		if time.Since(c.times[url]) < c.ttl {
			return result
		}
		delete(c.data, url)
		delete(c.times, url)
	}
	return nil
}

func (c *SimpleCache) Set(url string, result *models.WebAnalysingResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[url] = result
	c.times[url] = time.Now()
}

package github

import (
	"sync"
	"time"
)

const cacheTTL = 120 * time.Minute

type cacheEntry struct {
	data      interface{}
	createdAt time.Time
}

type CacheManager struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func NewCacheManager() *CacheManager {
	return &CacheManager{entries: make(map[string]cacheEntry)}
}

func (cm *CacheManager) Get(key string) (interface{}, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	e, ok := cm.entries[key]
	if !ok || time.Since(e.createdAt) > cacheTTL {
		return nil, false
	}
	return e.data, true
}

func (cm *CacheManager) Set(key string, data interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.entries[key] = cacheEntry{data: data, createdAt: time.Now()}
}

func (cm *CacheManager) CleanExpired() int {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	now := time.Now()
	cleaned := 0
	for k, e := range cm.entries {
		if now.Sub(e.createdAt) > cacheTTL {
			delete(cm.entries, k)
			cleaned++
		}
	}
	return cleaned
}

func (cm *CacheManager) Len() int {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	return len(cm.entries)
}

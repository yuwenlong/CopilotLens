package github

import (
	"sync"
	"time"
)

const (
	cacheTTL     = 120 * time.Minute
	maxCacheSize = 5000
)

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
	e, ok := cm.entries[key]
	cm.mu.RUnlock()
	if !ok || time.Since(e.createdAt) > cacheTTL {
		if ok {
			cm.mu.Lock()
			delete(cm.entries, key)
			cm.mu.Unlock()
		}
		return nil, false
	}
	return e.data, true
}

func (cm *CacheManager) Set(key string, data interface{}) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if len(cm.entries) >= maxCacheSize {
		var oldestKey string
		var oldestTime time.Time
		for k, e := range cm.entries {
			if oldestKey == "" || e.createdAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = e.createdAt
			}
		}
		if oldestKey != "" {
			delete(cm.entries, oldestKey)
		}
	}
	cm.entries[key] = cacheEntry{data: data, createdAt: time.Now()}
}

// Touch 刷新未过期条目的创建时间，实现续命；已过期或不存在则不处理。
func (cm *CacheManager) Touch(key string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	if e, ok := cm.entries[key]; ok && time.Since(e.createdAt) <= cacheTTL {
		cm.entries[key] = cacheEntry{data: e.data, createdAt: time.Now()}
	}
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

package tasks

import (
	"log"
	"time"

	"copilotlens/internal/github"
)

type CacheCleanTask struct {
	ticker *time.Ticker
	cache  *github.CacheManager
}

func NewCacheCleanTask(cache *github.CacheManager) *CacheCleanTask {
	return &CacheCleanTask{cache: cache}
}

func (t *CacheCleanTask) Run() {
	t.ticker = time.NewTicker(10 * time.Minute)
	go func() {
		n := t.cache.CleanExpired()
		log.Printf("缓存清理完成，清除 %d 条过期条目，剩余 %d 条", n, t.cache.Len())
		for range t.ticker.C {
			n := t.cache.CleanExpired()
			if n > 0 {
				log.Printf("缓存清理完成，清除 %d 条过期条目，剩余 %d 条", n, t.cache.Len())
			}
		}
	}()
}

func (t *CacheCleanTask) Stop() {
	t.ticker.Stop()
}

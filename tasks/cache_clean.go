package tasks

import (
	"log"
	"time"

	"copilotlens/internal/github"
)

type CacheCleanTask struct {
	ticker *time.Ticker
	cache  *github.CacheManager
	done   chan struct{}
}

func NewCacheCleanTask(cache *github.CacheManager) *CacheCleanTask {
	return &CacheCleanTask{cache: cache, done: make(chan struct{})}
}

func (t *CacheCleanTask) Run() {
	t.ticker = time.NewTicker(10 * time.Minute)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("CacheCleanTask panic: %v", r)
			}
		}()
		n := t.cache.CleanExpired()
		log.Printf("缓存清理完成，清除 %d 条过期条目，剩余 %d 条", n, t.cache.Len())
		for {
			select {
			case <-t.ticker.C:
				n := t.cache.CleanExpired()
				if n > 0 {
					log.Printf("缓存清理完成，清除 %d 条过期条目，剩余 %d 条", n, t.cache.Len())
				}
			case <-t.done:
				return
			}
		}
	}()
}

func (t *CacheCleanTask) Stop() {
	close(t.done)
	t.ticker.Stop()
}

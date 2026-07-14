package tasks

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"copilotlens/internal/github"
)

type CacheCleanTask struct {
	ticker       *time.Ticker
	cache        *github.CacheManager
	lastLogClean string
	done         chan struct{}
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
				today := time.Now().Format("2006-01-02")
				if today != t.lastLogClean {
					t.lastLogClean = today
					if r := t.cleanOldLogs("logs", 7); r > 0 {
						log.Printf("清理旧日志，删除 %d 个文件", r)
					}
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

func (t *CacheCleanTask) cleanOldLogs(dir string, maxDays int) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	cutoff := time.Now().AddDate(0, 0, -maxDays)
	removed := 0
	for _, e := range entries {
		if e.IsDir() || !strings.HasPrefix(e.Name(), "copilotlens-") || !strings.HasSuffix(e.Name(), ".log") {
			continue
		}
		dateStr := strings.TrimPrefix(e.Name(), "copilotlens-")
		dateStr = strings.TrimSuffix(dateStr, ".log")
		t, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			continue
		}
		if t.Before(cutoff) {
			os.Remove(filepath.Join(dir, e.Name()))
			removed++
		}
	}
	return removed
}

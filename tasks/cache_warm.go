package tasks

import (
	"log"
	"sync"
	"time"

	"copilotlens/internal/github"
)

const warmConcurrency = 20

type CacheWarmTask struct {
	ticker *time.Ticker
	client *github.Client
}

func NewCacheWarmTask(client *github.Client) *CacheWarmTask {
	return &CacheWarmTask{client: client}
}

func (t *CacheWarmTask) Run() {
	go func() {
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		time.Sleep(time.Until(nextHour))

		t.Warm()

		t.ticker = time.NewTicker(1 * time.Hour)
		for range t.ticker.C {
			t.Warm()
		}
	}()
}

func (t *CacheWarmTask) Stop() {
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *CacheWarmTask) Warm() {
	log.Println("开始缓存预热...")
	start := time.Now()
	now := time.Now()
	year, month, day := now.Year(), int(now.Month()), now.Day()

	users, err := t.client.GetOrgMembers()
	if err != nil {
		log.Printf("预热获取用户列表失败: %v", err)
		return
	}

	sem := make(chan struct{}, warmConcurrency)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()
		if _, err := t.client.GetMonthlyUsage(year, month); err != nil {
			log.Printf("预热月度用量失败: %v", err)
		}
	}()

	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	for d := startOfMonth; !d.After(now); d = d.AddDate(0, 0, 1) {
		yy, mm, dd := d.Year(), int(d.Month()), d.Day()
		wg.Add(1)
		go func(y, m, d int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetDailyUsage(y, m, d); err != nil {
				log.Printf("预热 %04d-%02d-%02d 每日用量失败: %v", y, m, d, err)
			}
		}(yy, mm, dd)
	}

	for _, u := range users {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetUserMonthlyUsage(username, year, month); err != nil {
				log.Printf("预热用户 %s 月用量失败: %v", username, err)
			}
		}(u)
	}

	for _, u := range users {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetUserDailyUsage(username, year, month, day); err != nil {
				log.Printf("预热用户 %s 当日用量失败: %v", username, err)
			}
		}(u)
	}

	wg.Wait()
	log.Printf("缓存预热完成，耗时 %v", time.Since(start))
}

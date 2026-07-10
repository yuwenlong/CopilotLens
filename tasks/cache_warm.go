package tasks

import (
	"log"
	"sync"
	"time"

	"copilotlens/internal/config"
	"copilotlens/internal/github"
	"copilotlens/internal/service"
)

const warmConcurrency = 20

type CacheWarmTask struct {
	mu     sync.Mutex
	ticker *time.Ticker
	client *github.Client
	conf   *config.AppConfig
	done   chan struct{}
}

func NewCacheWarmTask(client *github.Client) *CacheWarmTask {
	return &CacheWarmTask{
		client: client,
		conf:   config.Config(),
		done:   make(chan struct{}),
	}
}

func (t *CacheWarmTask) Run() {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("CacheWarmTask panic: %v", r)
			}
		}()
		now := time.Now()
		nextHour := now.Truncate(time.Hour).Add(time.Hour)
		select {
		case <-time.After(time.Until(nextHour)):
		case <-t.done:
			return
		}

		t.Warm()

		t.mu.Lock()
		t.ticker = time.NewTicker(1 * time.Hour)
		t.mu.Unlock()
		for {
			select {
			case <-t.ticker.C:
				t.Warm()
			case <-t.done:
				return
			}
		}
	}()
}

func (t *CacheWarmTask) Stop() {
	close(t.done)
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.ticker != nil {
		t.ticker.Stop()
	}
}

func (t *CacheWarmTask) Warm() {
	log.Println("开始缓存预热...")
	start := time.Now()
	now := time.Now().UTC()

	t.client.CheckRateLimit()

	if t.client.ShouldThrottle(100) {
		log.Printf("Rate limit 剩余 %d，跳过本次预热", t.client.RateLimitRemaining())
		return
	}

	users, err := t.client.GetOrgMembers()
	if err != nil {
		log.Printf("预热获取用户列表失败: %v", err)
		return
	}

	year, month := now.Year(), int(now.Month())
	t.warmMonth(year, month, users)

	log.Printf("缓存预热完成，耗时 %v", time.Since(start))
}

func (t *CacheWarmTask) warmMonth(year, month int, users []string) {
	now := time.Now().UTC()
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	var lastDay int
	if year == now.Year() && month == int(now.Month()) {
		lastDay = now.Day()
	} else {
		lastDay = endOfMonth.Day()
	}

	log.Printf("预热 %04d-%02d，日期范围 1-%d", year, month, lastDay)

	sem := make(chan struct{}, warmConcurrency)
	var wg sync.WaitGroup

	// monthly
	wg.Add(1)
	go func(y, m int) {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()
		if _, err := t.client.GetMonthlyUsage(y, m); err != nil {
			log.Printf("预热月度用量 %04d-%02d 失败: %v", y, m, err)
		}
	}(year, month)

	// daily full month
	for day := 1; day <= lastDay; day++ {
		wg.Add(1)
		go func(y, m, d int) {
			defer wg.Done()
			if t.client.ShouldThrottle(50) {
				return
			}
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetDailyUsage(y, m, d); err != nil {
				log.Printf("预热每日用量 %04d-%02d-%02d 失败: %v", y, m, d, err)
			}
		}(year, month, day)
	}

	// user monthly
	for _, username := range users {
		wg.Add(1)
		go func(username string, y, m int) {
			defer wg.Done()
			if t.client.ShouldThrottle(50) {
				return
			}
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetUserMonthlyUsage(username, y, m); err != nil {
				log.Printf("预热用户 %s 月用量 %04d-%02d 失败: %v", username, y, m, err)
			}
		}(username, year, month)
	}

	// user daily full month
	for _, username := range users {
		for day := 1; day <= lastDay; day++ {
			wg.Add(1)
			go func(username string, y, m, d int) {
				defer wg.Done()
				if t.client.ShouldThrottle(50) {
					log.Printf("Rate limit 剩余 %d，跳过用户 %s %04d-%02d-%02d 预热", t.client.RateLimitRemaining(), username, y, m, d)
					return
				}
				sem <- struct{}{}
				defer func() { <-sem }()
				if _, err := t.client.GetUserDailyUsage(username, y, m, d); err != nil {
					log.Printf("预热用户 %s 日用量 %04d-%02d-%02d 失败: %v", username, y, m, d, err)
				}
			}(username, year, month, day)
		}
	}

	wg.Wait()

	svc := service.NewUsageService(t.client, "data")
	svc.GetMonthlyTotal(year, month)
	svc.GetMonthlyUser(year, month)
	svc.GetMonthlyModel(year, month)
	for day := 1; day <= lastDay; day++ {
		svc.GetDailyUser(year, month, day)
		svc.GetDailyModel(year, month, day)
	}

	log.Printf("%04d-%02d 预热完成", year, month)
}

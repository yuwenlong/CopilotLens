package tasks

import (
	"log"
	"sync"
	"time"

	"copilotlens/internal/config"
	"copilotlens/internal/github"
	"copilotlens/internal/service"
)

const (
	warmConcurrency    = 20
	warmRetryIntervals = 3
)

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

		if !t.Warm() {
			t.retryWarm(0)
		}

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

func (t *CacheWarmTask) Warm() bool {
	log.Println("开始缓存预热...")
	start := time.Now()
	now := time.Now().UTC()

	t.client.CheckRateLimit()

	if t.client.ShouldThrottle(100) {
		log.Printf("Rate limit 剩余 %d，跳过本次预热", t.client.RateLimitRemaining())
		return false
	}

	users, err := t.client.GetOrgMembers()
	if err != nil {
		log.Printf("预热获取用户列表失败: %v", err)
		return false
	}

	year, month := now.Year(), int(now.Month())
	today := now.Day()
	t.warmMonth(year, month, today, users)

	log.Printf("缓存预热完成，耗时 %v", time.Since(start))
	return true
}

func (t *CacheWarmTask) retryWarm(attempt int) {
	intervals := []time.Duration{5 * time.Minute, 15 * time.Minute, 30 * time.Minute}
	if attempt >= warmRetryIntervals {
		log.Printf("预热重试次数已用尽（%d 次），放弃", attempt)
		return
	}
	wait := intervals[attempt]
	log.Printf("预热将在 %v 后进行第 %d 次重试", wait, attempt+1)
	time.AfterFunc(wait, func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("CacheWarmTask retry panic: %v", r)
			}
		}()
		if !t.Warm() {
			t.retryWarm(attempt + 1)
		}
	})
}

func (t *CacheWarmTask) warmMonth(year, month, today int, users []string) {
	log.Printf("预热 %04d-%02d，日期范围 1-%d，用户数 %d", year, month, today, len(users))

	sem := make(chan struct{}, warmConcurrency)
	var wg sync.WaitGroup

	// Phase 1: 原始数据预热

	// monthly usage (1 次)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sem <- struct{}{}
		defer func() { <-sem }()
		if _, err := t.client.GetMonthlyUsage(year, month); err != nil {
			log.Printf("预热月度用量 %04d-%02d 失败: %v", year, month, err)
		}
	}()

	// daily usage × today 天 (最多 31 次)
	for day := 1; day <= today; day++ {
		wg.Add(1)
		go func(d int) {
			defer wg.Done()
			if t.client.ShouldThrottle(50) {
				return
			}
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetDailyUsage(year, month, d); err != nil {
				log.Printf("预热每日用量 %04d-%02d-%02d 失败: %v", year, month, d, err)
			}
		}(day)
	}

	// user monthly × len(users) (~40 次)
	for _, username := range users {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			if t.client.ShouldThrottle(50) {
				return
			}
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetUserMonthlyUsage(u, year, month); err != nil {
				log.Printf("预热用户 %s 月用量 %04d-%02d 失败: %v", u, year, month, err)
			}
		}(username)
	}

	// user daily × len(users) × 仅今天 (~40 次)
	for _, username := range users {
		wg.Add(1)
		go func(u string) {
			defer wg.Done()
			if t.client.ShouldThrottle(50) {
				return
			}
			sem <- struct{}{}
			defer func() { <-sem }()
			if _, err := t.client.GetUserDailyUsage(u, year, month, today); err != nil {
				log.Printf("预热用户 %s 日用量 %04d-%02d-%02d 失败: %v", u, year, month, today, err)
			}
		}(username)
	}

	wg.Wait()

	// Phase 2: 聚合缓存构建
	svc := service.NewUsageService(t.client, "data")

	if _, err := svc.GetMonthlyTotal(year, month); err != nil {
		log.Printf("聚合 monthly_total 失败: %v", err)
	}
	if _, err := svc.GetMonthlyUser(year, month); err != nil {
		log.Printf("聚合 monthly_user 失败: %v", err)
	}
	if _, err := svc.GetMonthlyModel(year, month); err != nil {
		log.Printf("聚合 monthly_model 失败: %v", err)
	}
	if _, err := svc.GetDailyUser(year, month, today); err != nil {
		log.Printf("聚合 daily_user 失败: %v", err)
	}
	if _, err := svc.GetDailyModel(year, month, today); err != nil {
		log.Printf("聚合 daily_model 失败: %v", err)
	}

	log.Printf("%04d-%02d 预热完成", year, month)
}

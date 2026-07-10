package github

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	apiBase    = "https://api.github.com"
	apiVersion = "2022-11-28"
)

type BillingResponse struct {
	TimePeriod   TimePeriod  `json:"timePeriod"`
	Organization string      `json:"organization"`
	User         string      `json:"user,omitempty"`
	UsageItems   []UsageItem `json:"usageItems"`
}

type TimePeriod struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day,omitempty"`
}

type UsageItem struct {
	Product          string  `json:"product"`
	SKU              string  `json:"sku"`
	Model            string  `json:"model"`
	UnitType         string  `json:"unitType"`
	PricePerUnit     float64 `json:"pricePerUnit"`
	GrossQuantity    float64 `json:"grossQuantity"`
	GrossAmount      float64 `json:"grossAmount"`
	DiscountQuantity float64 `json:"discountQuantity"`
	DiscountAmount   float64 `json:"discountAmount"`
	NetQuantity      float64 `json:"netQuantity"`
	NetAmount        float64 `json:"netAmount"`
}

type OrgMember struct {
	Login string `json:"login"`
}

type Client struct {
	token      string
	org        string
	httpClient *http.Client
	Cache      *CacheManager
	sf         singleflight.Group

	rateMu        sync.Mutex
	rateRemaining int
	rateReset     int64
}

func NewClient(token, org string) *Client {
	cache := NewCacheManager()
	return &Client{
		token:      token,
		org:        org,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		Cache:      cache,
	}
}

func (c *Client) doRequest(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", apiVersion)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	c.rateMu.Lock()
	if remaining := resp.Header.Get("X-RateLimit-Remaining"); remaining != "" {
		if v, err := strconv.Atoi(remaining); err == nil {
			c.rateRemaining = v
		}
	}
	if reset := resp.Header.Get("X-RateLimit-Reset"); reset != "" {
		if v, err := strconv.ParseInt(reset, 10, 64); err == nil {
			c.rateReset = v
		}
	}
	c.rateMu.Unlock()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) buildURL(path string, params map[string]string) string {
	u := apiBase + path
	first := true
	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if first {
			u += "?" + k + "=" + url.QueryEscape(params[k])
			first = false
		} else {
			u += "&" + k + "=" + url.QueryEscape(params[k])
		}
	}
	return u
}

func (c *Client) GetMonthlyUsage(year, month int) ([]UsageItem, error) {
	key := fmt.Sprintf("monthly:%d:%d", year, month)
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]UsageItem), nil
	}

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.([]UsageItem), nil
		}
		params := map[string]string{
			"year":  fmt.Sprintf("%d", year),
			"month": fmt.Sprintf("%d", month),
		}
		url := c.buildURL(fmt.Sprintf("/organizations/%s/settings/billing/ai_credit/usage", c.org), params)
		body, err := c.doRequest(url)
		if err != nil {
			return nil, err
		}
		var resp BillingResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		c.Cache.Set(key, resp.UsageItems)
		return resp.UsageItems, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]UsageItem), nil
}

func (c *Client) GetUserMonthlyUsage(username string, year, month int) ([]UsageItem, error) {
	key := fmt.Sprintf("user_monthly:%s:%d:%d", username, year, month)
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]UsageItem), nil
	}

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.([]UsageItem), nil
		}
		params := map[string]string{
			"year":  fmt.Sprintf("%d", year),
			"month": fmt.Sprintf("%d", month),
			"user":  username,
		}
		url := c.buildURL(fmt.Sprintf("/organizations/%s/settings/billing/ai_credit/usage", c.org), params)
		body, err := c.doRequest(url)
		if err != nil {
			return nil, err
		}
		var resp BillingResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		c.Cache.Set(key, resp.UsageItems)
		return resp.UsageItems, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]UsageItem), nil
}

func (c *Client) GetDailyUsage(year, month, day int) ([]UsageItem, error) {
	key := fmt.Sprintf("daily:%d:%d:%d", year, month, day)
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]UsageItem), nil
	}

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.([]UsageItem), nil
		}
		params := map[string]string{
			"year":  fmt.Sprintf("%d", year),
			"month": fmt.Sprintf("%d", month),
			"day":   fmt.Sprintf("%d", day),
		}
		url := c.buildURL(fmt.Sprintf("/organizations/%s/settings/billing/ai_credit/usage", c.org), params)
		body, err := c.doRequest(url)
		if err != nil {
			return nil, err
		}
		var resp BillingResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		c.Cache.Set(key, resp.UsageItems)
		return resp.UsageItems, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]UsageItem), nil
}

func (c *Client) GetUserDailyUsage(username string, year, month, day int) ([]UsageItem, error) {
	key := fmt.Sprintf("user_daily:%s:%d:%d:%d", username, year, month, day)
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]UsageItem), nil
	}

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.([]UsageItem), nil
		}
		params := map[string]string{
			"year":  fmt.Sprintf("%d", year),
			"month": fmt.Sprintf("%d", month),
			"day":   fmt.Sprintf("%d", day),
			"user":  username,
		}
		url := c.buildURL(fmt.Sprintf("/organizations/%s/settings/billing/ai_credit/usage", c.org), params)
		body, err := c.doRequest(url)
		if err != nil {
			return nil, err
		}
		var resp BillingResponse
		if err := json.Unmarshal(body, &resp); err != nil {
			return nil, err
		}
		c.Cache.Set(key, resp.UsageItems)
		return resp.UsageItems, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]UsageItem), nil
}

func (c *Client) GetOrgMembers() ([]string, error) {
	key := "org_members"
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]string), nil
	}

	v, err, _ := c.sf.Do(key, func() (interface{}, error) {
		if cached, ok := c.Cache.Get(key); ok {
			return cached.([]string), nil
		}

		var allMembers []string
		page := 1
		perPage := 100

		for {
			if c.ShouldThrottle(50) {
				log.Printf("GetOrgMembers: Rate limit 剩余 %d，停止分页", c.RateLimitRemaining())
				return nil, fmt.Errorf("rate limit exhausted, partial members fetched (page %d)", page)
			}
			url := fmt.Sprintf("%s/orgs/%s/members?per_page=%d&page=%d", apiBase, c.org, perPage, page)
			body, err := c.doRequest(url)
			if err != nil {
				return nil, err
			}

			var members []OrgMember
			if err := json.Unmarshal(body, &members); err != nil {
				return nil, err
			}
			if len(members) == 0 {
				break
			}

			for _, m := range members {
				allMembers = append(allMembers, m.Login)
			}

			if len(members) < perPage {
				break
			}
			page++
		}

		c.Cache.Set(key, allMembers)
		return allMembers, nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]string), nil
}

func (c *Client) RateLimitRemaining() int {
	c.rateMu.Lock()
	defer c.rateMu.Unlock()
	return c.rateRemaining
}

func (c *Client) ShouldThrottle(threshold int) bool {
	c.rateMu.Lock()
	remaining := c.rateRemaining
	reset := c.rateReset
	c.rateMu.Unlock()
	// 尚未收到任何响应（rateReset == 0），不 throttle，让调用方自行决定
	if reset == 0 {
		return false
	}
	return remaining < threshold
}

func (c *Client) CheckRateLimit() {
	if _, err := c.doRequest(apiBase); err != nil {
		log.Printf("CheckRateLimit 失败: %v", err)
	}
}

func (c *Client) WaitForRateLimit() {
	c.rateMu.Lock()
	resetTime := time.Unix(c.rateReset, 0)
	c.rateMu.Unlock()
	wait := time.Until(resetTime) + 5*time.Second
	const maxWait = 5 * time.Minute
	if wait > maxWait {
		wait = maxWait
	}
	if wait > 0 {
		log.Printf("Rate limit 剩余 %d，等待 %v", c.RateLimitRemaining(), wait)
		time.Sleep(wait)
	}
}

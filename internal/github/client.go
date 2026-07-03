package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
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
}

func NewClient(token, org string) *Client {
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
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

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API error %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}

func (c *Client) buildURL(path string, params map[string]string) string {
	url := apiBase + path
	first := true
	for k, v := range params {
		if first {
			url += "?" + k + "=" + v
			first = false
		} else {
			url += "&" + k + "=" + v
		}
	}
	return url
}

func (c *Client) GetMonthlyUsage(year, month int) ([]UsageItem, error) {
	key := fmt.Sprintf("monthly:%d:%d", year, month)
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
}

func (c *Client) GetUserMonthlyUsage(username string, year, month int) ([]UsageItem, error) {
	key := fmt.Sprintf("user_monthly:%s:%d:%d", username, year, month)
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
}

func (c *Client) GetDailyUsage(year, month, day int) ([]UsageItem, error) {
	key := fmt.Sprintf("daily:%d:%d:%d", year, month, day)
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
}

func (c *Client) GetUserDailyUsage(username string, year, month, day int) ([]UsageItem, error) {
	key := fmt.Sprintf("user_daily:%s:%d:%d:%d", username, year, month, day)
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
}

func (c *Client) GetOrgMembers() ([]string, error) {
	key := "org_members"
	if cached, ok := c.Cache.Get(key); ok {
		return cached.([]string), nil
	}

	var allMembers []string
	page := 1
	perPage := 100

	for {
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
}

func (c *Client) GetDailyUsageConcurrent(year, month int, days []int) map[int][]UsageItem {
	result := make(map[int][]UsageItem)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, day := range days {
		wg.Add(1)
		go func(d int) {
			defer wg.Done()
			items, err := c.GetDailyUsage(year, month, d)
			if err != nil {
				return
			}
			mu.Lock()
			result[d] = items
			mu.Unlock()
		}(day)
	}
	wg.Wait()
	return result
}

func (c *Client) GetUserMonthlyUsageConcurrent(users []string, year, month int) map[string][]UsageItem {
	result := make(map[string][]UsageItem)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, u := range users {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			items, err := c.GetUserMonthlyUsage(username, year, month)
			if err != nil {
				return
			}
			mu.Lock()
			result[username] = items
			mu.Unlock()
		}(u)
	}
	wg.Wait()
	return result
}

func (c *Client) GetUserDailyUsageConcurrent(users []string, year, month, day int) map[string][]UsageItem {
	result := make(map[string][]UsageItem)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, u := range users {
		wg.Add(1)
		go func(username string) {
			defer wg.Done()
			items, err := c.GetUserDailyUsage(username, year, month, day)
			if err != nil {
				return
			}
			mu.Lock()
			result[username] = items
			mu.Unlock()
		}(u)
	}
	wg.Wait()
	return result
}

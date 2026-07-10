package service

import (
	"fmt"
	"log"
	"sort"
	"time"

	"copilotlens/domain/dto"
	"copilotlens/internal/client"
	"copilotlens/internal/github"
	"golang.org/x/sync/singleflight"
)

const (
	respMonthlyTotalKey = "resp:monthly_total:%04d-%02d"
	respMonthlyUserKey  = "resp:monthly_user:%04d-%02d"
	respDailyUserKey    = "resp:daily_user:%04d-%02d-%02d"
	respMonthlyModelKey = "resp:monthly_model:%04d-%02d"
	respDailyModelKey   = "resp:daily_model:%04d-%02d-%02d"
)

type UsageService struct {
	client  *github.Client
	DataDir string
	sf      singleflight.Group
}

func NewUsageService(client *github.Client, dataDir string) *UsageService {
	return &UsageService{client: client, DataDir: dataDir}
}

func MonthlyTotalCacheKey(year, month int) string {
	return fmt.Sprintf(respMonthlyTotalKey, year, month)
}

func MonthlyUserCacheKey(year, month int) string {
	return fmt.Sprintf(respMonthlyUserKey, year, month)
}

func DailyUserCacheKey(year, month, day int) string {
	return fmt.Sprintf(respDailyUserKey, year, month, day)
}

func MonthlyModelCacheKey(year, month int) string {
	return fmt.Sprintf(respMonthlyModelKey, year, month)
}

func DailyModelCacheKey(year, month, day int) string {
	return fmt.Sprintf(respDailyModelKey, year, month, day)
}

func (s *UsageService) GetMonthlyTotal(year, month int) (*dto.MonthlyTotalResponse, error) {
	key := MonthlyTotalCacheKey(year, month)
	if v, ok := s.client.Cache.Get(key); ok {
		return v.(*dto.MonthlyTotalResponse), nil
	}
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		if v, ok := s.client.Cache.Get(key); ok {
			return v.(*dto.MonthlyTotalResponse), nil
		}
		resp, err := s.BuildMonthlyTotal(year, month)
		if err != nil {
			return nil, err
		}
		s.client.Cache.Set(key, resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*dto.MonthlyTotalResponse), nil
}

func (s *UsageService) BuildMonthlyTotal(year, month int) (*dto.MonthlyTotalResponse, error) {
	items, err := s.client.GetMonthlyUsage(year, month)
	if err != nil {
		return nil, err
	}

	var total float64
	for _, item := range items {
		total += item.GrossQuantity
	}

	daily := make([]dto.DailyTotal, 0)
	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	end := start.AddDate(0, 1, 0)
	now := time.Now()
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		if d.After(now) {
			break
		}
		items, err := s.client.GetDailyUsage(d.Year(), int(d.Month()), d.Day())
		var dayTotal float64
		if err != nil {
			log.Printf("BuildMonthlyTotal: 获取每日用量 %04d-%02d-%02d 失败: %v", d.Year(), int(d.Month()), d.Day(), err)
		} else {
			for _, item := range items {
				dayTotal += item.GrossQuantity
			}
		}
		daily = append(daily, dto.DailyTotal{Date: d.Format("2006-01-02"), Total: client.Round2(dayTotal)})
	}
	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Date > daily[j].Date
	})

	return &dto.MonthlyTotalResponse{
		Month: fmt.Sprintf("%04d-%02d", year, month),
		Total: client.Round2(total),
		Daily: daily,
	}, nil
}

func (s *UsageService) GetMonthlyUser(year, month int) (*dto.MonthlyUserResponse, error) {
	key := MonthlyUserCacheKey(year, month)
	if v, ok := s.client.Cache.Get(key); ok {
		return v.(*dto.MonthlyUserResponse), nil
	}
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		if v, ok := s.client.Cache.Get(key); ok {
			return v.(*dto.MonthlyUserResponse), nil
		}
		resp, err := s.BuildMonthlyUser(year, month)
		if err != nil {
			return nil, err
		}
		s.client.Cache.Set(key, resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*dto.MonthlyUserResponse), nil
}

func (s *UsageService) BuildMonthlyUser(year, month int) (*dto.MonthlyUserResponse, error) {
	usernameMap := client.LoadUsernameMap(s.DataDir)
	users := client.LoadUsersFromClient(s.client)

	allUsers := make([]dto.UserUsage, 0)
	for _, username := range users {
		items, err := s.client.GetUserMonthlyUsage(username, year, month)
		if err != nil {
			log.Printf("BuildMonthlyUser: 获取用户 %s 月用量失败: %v", username, err)
			continue
		}
		if len(items) == 0 {
			continue
		}
		allUsers = append(allUsers, itemsToUserUsage(username, usernameMap, items))
	}

	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].Total > allUsers[j].Total
	})

	return &dto.MonthlyUserResponse{
		Month: fmt.Sprintf("%04d-%02d", year, month),
		Users: allUsers,
	}, nil
}

func (s *UsageService) GetDailyUser(year, month, day int) (*dto.MonthlyUserResponse, error) {
	key := DailyUserCacheKey(year, month, day)
	if v, ok := s.client.Cache.Get(key); ok {
		return v.(*dto.MonthlyUserResponse), nil
	}
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		if v, ok := s.client.Cache.Get(key); ok {
			return v.(*dto.MonthlyUserResponse), nil
		}
		resp, err := s.BuildDailyUser(year, month, day)
		if err != nil {
			return nil, err
		}
		s.client.Cache.Set(key, resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*dto.MonthlyUserResponse), nil
}

func (s *UsageService) BuildDailyUser(year, month, day int) (*dto.MonthlyUserResponse, error) {
	usernameMap := client.LoadUsernameMap(s.DataDir)
	users := client.LoadUsersFromClient(s.client)

	allUsers := make([]dto.UserUsage, 0)
	for _, username := range users {
		items, err := s.client.GetUserDailyUsage(username, year, month, day)
		if err != nil {
			log.Printf("BuildDailyUser: 获取用户 %s 日用量 %04d-%02d-%02d 失败: %v", username, year, month, day, err)
			continue
		}
		if len(items) == 0 {
			continue
		}
		allUsers = append(allUsers, itemsToUserUsage(username, usernameMap, items))
	}

	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].Total > allUsers[j].Total
	})

	return &dto.MonthlyUserResponse{
		Month: fmt.Sprintf("%04d-%02d-%02d", year, month, day),
		Users: allUsers,
	}, nil
}

func (s *UsageService) GetMonthlyModel(year, month int) (*dto.MonthlyModelResponse, error) {
	key := MonthlyModelCacheKey(year, month)
	if v, ok := s.client.Cache.Get(key); ok {
		return v.(*dto.MonthlyModelResponse), nil
	}
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		if v, ok := s.client.Cache.Get(key); ok {
			return v.(*dto.MonthlyModelResponse), nil
		}
		resp, err := s.BuildMonthlyModel(year, month)
		if err != nil {
			return nil, err
		}
		s.client.Cache.Set(key, resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*dto.MonthlyModelResponse), nil
}

func (s *UsageService) BuildMonthlyModel(year, month int) (*dto.MonthlyModelResponse, error) {
	items, err := s.client.GetMonthlyUsage(year, month)
	if err != nil {
		return nil, err
	}
	models := aggregateModels(items)
	return &dto.MonthlyModelResponse{
		Month:  fmt.Sprintf("%04d-%02d", year, month),
		Models: models,
	}, nil
}

func (s *UsageService) GetDailyModel(year, month, day int) (*dto.MonthlyModelResponse, error) {
	key := DailyModelCacheKey(year, month, day)
	if v, ok := s.client.Cache.Get(key); ok {
		return v.(*dto.MonthlyModelResponse), nil
	}
	v, err, _ := s.sf.Do(key, func() (interface{}, error) {
		if v, ok := s.client.Cache.Get(key); ok {
			return v.(*dto.MonthlyModelResponse), nil
		}
		resp, err := s.BuildDailyModel(year, month, day)
		if err != nil {
			return nil, err
		}
		s.client.Cache.Set(key, resp)
		return resp, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*dto.MonthlyModelResponse), nil
}

func (s *UsageService) BuildDailyModel(year, month, day int) (*dto.MonthlyModelResponse, error) {
	items, err := s.client.GetDailyUsage(year, month, day)
	if err != nil {
		return nil, err
	}
	models := aggregateModels(items)
	return &dto.MonthlyModelResponse{
		Month:  fmt.Sprintf("%04d-%02d-%02d", year, month, day),
		Models: models,
	}, nil
}

func aggregateModels(items []github.UsageItem) []dto.ModelUsage {
	modelMap := make(map[string]*dto.ModelUsage)
	for _, item := range items {
		m, ok := modelMap[item.Model]
		if !ok {
			m = &dto.ModelUsage{Model: item.Model}
			modelMap[item.Model] = m
		}
		m.Total += item.GrossQuantity
		m.Cost += item.GrossAmount
	}

	models := make([]dto.ModelUsage, 0)
	for _, m := range modelMap {
		models = append(models, dto.ModelUsage{
			Model: m.Model,
			Total: client.Round2(m.Total),
			Cost:  client.Round2(m.Cost),
		})
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].Total > models[j].Total
	})
	return models
}

func itemsToUserUsage(username string, usernameMap map[string]string, items []github.UsageItem) dto.UserUsage {
	var total, cost float64
	modelMap := make(map[string]*dto.UserModel)
	for _, item := range items {
		total += item.GrossQuantity
		cost += item.GrossAmount
		m, ok := modelMap[item.Model]
		if !ok {
			m = &dto.UserModel{Model: item.Model}
			modelMap[item.Model] = m
		}
		m.Total += item.GrossQuantity
		m.Cost += item.GrossAmount
	}
	models := make([]dto.UserModel, 0)
	for _, m := range modelMap {
		models = append(models, dto.UserModel{
			Model: m.Model,
			Total: client.Round2(m.Total),
			Cost:  client.Round2(m.Cost),
		})
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].Total > models[j].Total
	})
	name := username
	if v, ok := usernameMap[username]; ok && v != "" {
		name = v
	}
	return dto.UserUsage{
		Username: username,
		Name:     name,
		Total:    client.Round2(total),
		Cost:     client.Round2(cost),
		Models:   models,
	}
}

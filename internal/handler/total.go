package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"copilotlens/domain/dto"
	"copilotlens/internal/client"
	"copilotlens/internal/github"
)

func (h *Handler) MonthlyTotal(c *gin.Context) {
	month := c.Query("month")
	if len(month) != 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month format must be YYYY-MM"})
		return
	}

	parts := strings.Split(month, "-")
	year, _ := strconv.Atoi(parts[0])
	mon, _ := strconv.Atoi(parts[1])

	items, err := client.GitHubClient.GetMonthlyUsage(year, mon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var total float64
	for _, item := range items {
		total += item.GrossQuantity
	}

	var daily []dto.DailyTotal
	start := time.Date(year, time.Month(mon), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)
	now := time.Now()
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		if d.After(now) {
			break
		}
		items, err := client.GitHubClient.GetDailyUsage(d.Year(), int(d.Month()), d.Day())
		var dayTotal float64
		if err == nil {
			for _, item := range items {
				dayTotal += item.GrossQuantity
			}
		}
		daily = append(daily, dto.DailyTotal{Date: d.Format("2006-01-02"), Total: client.Round2(dayTotal)})
	}
	sort.Slice(daily, func(i, j int) bool {
		return daily[i].Date > daily[j].Date
	})

	c.JSON(http.StatusOK, dto.MonthlyTotalResponse{
		Month: month,
		Total: client.Round2(total),
		Daily: daily,
	})
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
	var models []dto.UserModel
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

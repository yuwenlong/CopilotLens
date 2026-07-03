package handler

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"copilotlens/domain/dto"
	"copilotlens/internal/client"
)

func (h *Handler) MonthlyModel(c *gin.Context) {
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

	var models []dto.ModelUsage
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

	c.JSON(http.StatusOK, dto.MonthlyModelResponse{Month: month, Models: models})
}

func (h *Handler) DailyModel(c *gin.Context) {
	date := c.Query("date")
	if len(date) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date format must be YYYY-MM-DD"})
		return
	}

	parts := strings.Split(date, "-")
	year, _ := strconv.Atoi(parts[0])
	mon, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	items, err := client.GitHubClient.GetDailyUsage(year, mon, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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

	var models []dto.ModelUsage
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

	c.JSON(http.StatusOK, dto.MonthlyModelResponse{Month: date, Models: models})
}

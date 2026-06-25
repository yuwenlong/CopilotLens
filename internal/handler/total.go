package handler

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"copilotlens/domain/dto"
	"copilotlens/internal/client"
)

func (h *Handler) MonthlyTotal(c *gin.Context) {
	month := c.Query("month")
	if len(month) != 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month format must be YYYY-MM"})
		return
	}
	records := client.LoadCopilotCSV(h.DataDir, month)
	var total float64
	dailyMap := make(map[string]float64)
	for _, r := range records {
		total += r.AICQuantity
		dailyMap[r.Date] += r.AICQuantity
	}
	var daily []dto.DailyTotal
	for date, t := range dailyMap {
		daily = append(daily, dto.DailyTotal{Date: date, Total: client.Round2(t)})
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

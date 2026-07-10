package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MonthlyUser(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	month := c.Query("month")
	if len(month) != 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month format must be YYYY-MM"})
		return
	}

	parts := strings.Split(month, "-")
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	mon, err := strconv.Atoi(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}

	resp, err := h.usage.GetMonthlyUser(year, mon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *Handler) DailyUser(c *gin.Context) {
	c.Header("Cache-Control", "no-cache")
	date := c.Query("date")
	if len(date) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date format must be YYYY-MM-DD"})
		return
	}

	parts := strings.Split(date, "-")
	year, err := strconv.Atoi(parts[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year"})
		return
	}
	mon, err := strconv.Atoi(parts[1])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month"})
		return
	}
	day, err := strconv.Atoi(parts[2])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid day"})
		return
	}

	resp, err := h.usage.GetDailyUser(year, mon, day)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) MonthlyTotal(c *gin.Context) {
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

	resp, err := h.usage.GetMonthlyTotal(year, mon)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

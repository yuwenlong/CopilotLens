package handler

import (
	"net/http"

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
	for _, r := range records {
		total += r.AICQuantity
	}
	c.JSON(http.StatusOK, dto.MonthlyTotalResponse{
		Month: month,
		Total: client.Round2(total),
	})
}

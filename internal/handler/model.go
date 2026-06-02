package handler

import (
	"net/http"
	"sort"

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
	records := client.LoadCopilotCSV(h.DataDir, month)

	type modelAgg struct {
		total float64
		cost  float64
	}
	modelAggMap := make(map[string]*modelAgg)

	for _, r := range records {
		ma, ok := modelAggMap[r.Model]
		if !ok {
			ma = &modelAgg{}
			modelAggMap[r.Model] = ma
		}
		ma.total += r.AICQuantity
		ma.cost += r.AICCost
	}

	var models []dto.ModelUsage
	for name, ma := range modelAggMap {
		models = append(models, dto.ModelUsage{
			Model: name,
			Total: client.Round2(ma.total),
			Cost:  client.Round2(ma.cost),
		})
	}
	sort.Slice(models, func(i, j int) bool {
		return models[i].Total > models[j].Total
	})

	c.JSON(http.StatusOK, dto.MonthlyModelResponse{Month: month, Models: models})
}

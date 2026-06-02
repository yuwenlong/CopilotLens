package handler

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"copilotlens/domain/dto"
	"copilotlens/internal/client"
)

func (h *Handler) MonthlyUser(c *gin.Context) {
	month := c.Query("month")
	if len(month) != 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month format must be YYYY-MM"})
		return
	}
	records := client.LoadCopilotCSV(h.DataDir, month)
	usernameMap := client.LoadUsernameMap(h.DataDir)

	type userModelAgg struct {
		models map[string]*dto.UserModel
		total  float64
		cost   float64
	}
	userAgg := make(map[string]*userModelAgg)

	for _, r := range records {
		ua, ok := userAgg[r.Username]
		if !ok {
			ua = &userModelAgg{models: make(map[string]*dto.UserModel)}
			userAgg[r.Username] = ua
		}
		ua.total += r.AICQuantity
		ua.cost += r.AICCost
		m, ok := ua.models[r.Model]
		if !ok {
			m = &dto.UserModel{Model: r.Model}
			ua.models[r.Model] = m
		}
		m.Total += r.AICQuantity
		m.Cost += r.AICCost
	}

	var users []dto.UserUsage
	for uname, ua := range userAgg {
		var models []dto.UserModel
		for _, m := range ua.models {
			models = append(models, dto.UserModel{
				Model: m.Model,
				Total: client.Round2(m.Total),
				Cost:  client.Round2(m.Cost),
			})
		}
		sort.Slice(models, func(i, j int) bool {
			return models[i].Total > models[j].Total
		})
		name := uname
		if v, ok := usernameMap[uname]; ok && v != "" {
			name = v
		}
		users = append(users, dto.UserUsage{
			Username: uname,
			Name:     name,
			Total:    client.Round2(ua.total),
			Cost:     client.Round2(ua.cost),
			Models:   models,
		})
	}
	sort.Slice(users, func(i, j int) bool {
		return users[i].Total > users[j].Total
	})

	c.JSON(http.StatusOK, dto.MonthlyUserResponse{Month: month, Users: users})
}

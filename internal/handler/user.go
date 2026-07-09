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

func (h *Handler) MonthlyUser(c *gin.Context) {
	month := c.Query("month")
	if len(month) != 7 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "month format must be YYYY-MM"})
		return
	}

	parts := strings.Split(month, "-")
	year, _ := strconv.Atoi(parts[0])
	mon, _ := strconv.Atoi(parts[1])

	usernameMap := client.LoadUsernameMap(h.DataDir)
	users := client.LoadUsers(h.DataDir)

	var allUsers []dto.UserUsage
	for _, username := range users {
		items, err := client.GitHubClient.GetUserMonthlyUsage(username, year, mon)
		if err != nil || len(items) == 0 {
			continue
		}
		allUsers = append(allUsers, itemsToUserUsage(username, usernameMap, items))
	}

	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].Total > allUsers[j].Total
	})

	c.JSON(http.StatusOK, dto.MonthlyUserResponse{Month: month, Users: allUsers})
}

func (h *Handler) DailyUser(c *gin.Context) {
	date := c.Query("date")
	if len(date) != 10 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date format must be YYYY-MM-DD"})
		return
	}

	parts := strings.Split(date, "-")
	year, _ := strconv.Atoi(parts[0])
	mon, _ := strconv.Atoi(parts[1])
	day, _ := strconv.Atoi(parts[2])

	usernameMap := client.LoadUsernameMap(h.DataDir)
	users := client.LoadUsers(h.DataDir)

	var allUsers []dto.UserUsage
	for _, username := range users {
		items, err := client.GitHubClient.GetUserDailyUsage(username, year, mon, day)
		if err != nil || len(items) == 0 {
			continue
		}
		allUsers = append(allUsers, itemsToUserUsage(username, usernameMap, items))
	}

	sort.Slice(allUsers, func(i, j int) bool {
		return allUsers[i].Total > allUsers[j].Total
	})

	c.JSON(http.StatusOK, dto.MonthlyUserResponse{Month: date, Users: allUsers})
}

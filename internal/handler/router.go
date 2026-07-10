package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/config"
	"copilotlens/internal/github"
	"copilotlens/internal/service"
)

var conf = config.Config()

type Handler struct {
	DataDir string
	cache   *github.CacheManager
	usage   *service.UsageService
}

func NewHandler(dataDir string, client *github.Client) *Handler {
	return &Handler{
		DataDir: dataDir,
		cache:   client.Cache,
		usage:   service.NewUsageService(client, dataDir),
	}
}

func IPWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Server.IsAccessAllowed(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ip not allowed"})
			return
		}
		c.Next()
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		h.renderPage(c, "index.html", "")
	})
	r.GET("/monthly-total", func(c *gin.Context) {
		month := c.Query("month")
		if month == "" {
			month = currentMonthStr()
		}
		h.renderPage(c, "monthly-total.html", monthlyTotalPageKey(month))
	})
	r.GET("/monthly-user", func(c *gin.Context) {
		month := c.Query("month")
		date := c.Query("date")
		if date == "" && month == "" {
			month = currentMonthStr()
		}
		h.renderPage(c, "monthly-user.html", monthlyUserPageKey(month, date))
	})
	r.GET("/monthly-model", func(c *gin.Context) {
		month := c.Query("month")
		date := c.Query("date")
		if date == "" && month == "" {
			month = currentMonthStr()
		}
		h.renderPage(c, "monthly-model.html", monthlyModelPageKey(month, date))
	})

	r.GET("/api/monthly-total", h.MonthlyTotal)
	r.GET("/api/monthly-user", h.MonthlyUser)
	r.GET("/api/monthly-model", h.MonthlyModel)
	r.GET("/api/daily-user", h.DailyUser)
	r.GET("/api/daily-model", h.DailyModel)
}

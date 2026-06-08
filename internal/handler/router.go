package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/conf"
)

type Handler struct {
	DataDir string
}

func IPWhitelist() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !conf.Cfg.Server.IsAccessAllowed(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "ip not allowed"})
			return
		}
		c.Next()
	}
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/monthly-total", func(c *gin.Context) {
		c.HTML(200, "monthly-total.html", nil)
	})
	r.GET("/monthly-user", func(c *gin.Context) {
		c.HTML(200, "monthly-user.html", nil)
	})
	r.GET("/monthly-model", func(c *gin.Context) {
		c.HTML(200, "monthly-model.html", nil)
	})

	r.GET("/api/monthly-total", h.MonthlyTotal)
	r.GET("/api/monthly-user", h.MonthlyUser)
	r.GET("/api/monthly-model", h.MonthlyModel)
	r.GET("/api/daily-user", h.DailyUser)
	r.GET("/api/daily-model", h.DailyModel)
}

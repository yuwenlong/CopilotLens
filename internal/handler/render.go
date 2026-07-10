package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"copilotlens/internal/service"
)

// pageInitialData 是注入 HTML 的直出数据结构。
type pageInitialData struct {
	MonthlyTotal interface{} `json:"monthly-total,omitempty"`
	MonthlyUser  interface{} `json:"monthly-user,omitempty"`
	DailyUser    interface{} `json:"daily-user,omitempty"`
	MonthlyModel interface{} `json:"monthly-model,omitempty"`
	DailyModel   interface{} `json:"daily-model,omitempty"`
}

// monthlyTotalPageKey 解析 month 参数并返回对应聚合缓存键，空字符串表示无效参数。
func monthlyTotalPageKey(month string) string {
	if len(month) != 7 {
		return ""
	}
	parts := strings.Split(month, "-")
	if len(parts) != 2 {
		return ""
	}
	year, err1 := strconv.Atoi(parts[0])
	mon, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		return ""
	}
	return service.MonthlyTotalCacheKey(year, mon)
}

// parseDateOrMonth 解析 YYYY-MM-DD 或 YYYY-MM 格式，返回 (year, month, day, ok)。
// day > 0 表示日期格式，day == 0 表示月份格式。
func parseDateOrMonth(date, month string) (year, mon, day int, ok bool) {
	if len(date) == 10 {
		parts := strings.Split(date, "-")
		if len(parts) == 3 {
			y, err1 := strconv.Atoi(parts[0])
			m, err2 := strconv.Atoi(parts[1])
			d, err3 := strconv.Atoi(parts[2])
			if err1 == nil && err2 == nil && err3 == nil {
				return y, m, d, true
			}
		}
	}
	if len(month) == 7 {
		parts := strings.Split(month, "-")
		if len(parts) == 2 {
			y, err1 := strconv.Atoi(parts[0])
			m, err2 := strconv.Atoi(parts[1])
			if err1 == nil && err2 == nil {
				return y, m, 0, true
			}
		}
	}
	return 0, 0, 0, false
}

func monthlyUserPageKey(month, date string) string {
	year, mon, day, ok := parseDateOrMonth(date, month)
	if !ok {
		return ""
	}
	if day > 0 {
		return service.DailyUserCacheKey(year, mon, day)
	}
		return service.MonthlyUserCacheKey(year, mon)
}

func monthlyModelPageKey(month, date string) string {
	year, mon, day, ok := parseDateOrMonth(date, month)
	if !ok {
		return ""
	}
	if day > 0 {
		return service.DailyModelCacheKey(year, mon, day)
	}
		return service.MonthlyModelCacheKey(year, mon)
}

func currentMonthStr() string {
	t := time.Now()
	return fmt.Sprintf("%04d-%02d", t.Year(), int(t.Month()))
}

func (h *Handler) renderPage(c *gin.Context, templateName, respKey string) {
	data := pageInitialData{}

	switch {
	case strings.HasPrefix(respKey, "resp:monthly_total:"):
		data.MonthlyTotal = h.cachedOrNull(respKey)
	case strings.HasPrefix(respKey, "resp:monthly_user:"):
		data.MonthlyUser = h.cachedOrNull(respKey)
	case strings.HasPrefix(respKey, "resp:daily_user:"):
		data.DailyUser = h.cachedOrNull(respKey)
	case strings.HasPrefix(respKey, "resp:monthly_model:"):
		data.MonthlyModel = h.cachedOrNull(respKey)
	case strings.HasPrefix(respKey, "resp:daily_model:"):
		data.DailyModel = h.cachedOrNull(respKey)
	}

	initialJSON, err := json.Marshal(data)
	if err != nil {
		initialJSON = []byte("{}")
	}

	html := InjectInitialData(string(initialJSON))

	c.HTML(http.StatusOK, templateName, gin.H{
		"InitialData": template.HTML(html),
	})
}

func (h *Handler) cachedOrNull(key string) interface{} {
	if cached, ok := h.cache.Get(key); ok {
		return cached
	}
	return nil
}



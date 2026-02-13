package handlers

import (
	"net/http"
	"time"

	"go-backend/internal/services"

	"github.com/gin-gonic/gin"
)

// AnalyticsHandler handles analytics endpoints
type AnalyticsHandler struct {
	service *services.AnalyticsService
}

// Constructor
func NewAnalyticsHandler(service *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) DailySummary(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	date = date.UTC()

	data, err := h.service.DailySummary(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *AnalyticsHandler) Trend(c *gin.Context) {
	fromStr := c.Query("from")
	toStr := c.Query("to")

	if fromStr == "" || toStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to dates are required"})
		return
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date format"})
		return
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date format"})
		return
	}

	from = from.UTC()
	to = to.UTC()

	data, err := h.service.Trend(from, to)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"labels": func() []string {
			var l []string
			for _, d := range data {
				l = append(l, d.Date.Format("2006-01-02"))
			}
			return l
		}(),
		"values": func() []int64 {
			var v []int64
			for _, d := range data {
				v = append(v, d.Count)
			}
			return v
		}(),
	})
}

func (h *AnalyticsHandler) Absentees(c *gin.Context) {
	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date is required"})
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format"})
		return
	}

	date = date.UTC()

	names, err := h.service.Absentees(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"absentees": names})
}

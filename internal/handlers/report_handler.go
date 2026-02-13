package handlers

import (
	"time"

	"github.com/gin-gonic/gin"

	"go-backend/internal/services"
)

type ReportHandler struct {
	reportService *services.ReportService
}

func NewReportHandler(reportService *services.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) ExportCSV(c *gin.Context) {
	from, _ := time.Parse("2006-01-02", c.Query("from"))
	to, _ := time.Parse("2006-01-02", c.Query("to"))

	c.Header("Content-Disposition", "attachment; filename=attendance.csv")
	c.Header("Content-Type", "text/csv")

	h.reportService.ExportCSV(c.Writer, from, to) // <- fixed field name
}

func (h *ReportHandler) ExportPDF(c *gin.Context) {
	from, _ := time.Parse("2006-01-02", c.Query("from"))
	to, _ := time.Parse("2006-01-02", c.Query("to"))

	c.Header("Content-Disposition", "attachment; filename=attendance.pdf")
	c.Header("Content-Type", "application/pdf")

	h.reportService.ExportPDF(c.Writer, from, to) // <- fixed field name
}

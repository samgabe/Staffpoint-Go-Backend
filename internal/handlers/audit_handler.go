package handlers

import (
	"encoding/csv"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"

	"go-backend/internal/services"
)

type AuditHandler struct {
	service services.AuditService
}

// Constructor
func NewAuditHandler(service services.AuditService) *AuditHandler {
	return &AuditHandler{service: service}
}

func (h *AuditHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	fromStr := c.Query("from")
	toStr := c.Query("to")
	userIDStr := c.Query("user_id")
	action := c.Query("action")

	var from, to *time.Time
	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err == nil {
			to = &t
		}
	}

	var userID *string
	if userIDStr != "" {
		userID = &userIDStr
	}

	logs, err := h.service.ListFiltered(limit, from, to, userID, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}

func (h *AuditHandler) ExportCSV(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "500"))
	from, to, userID, action := parseAuditFilters(c)

	logs, err := h.service.ListFiltered(limit, from, to, userID, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=audit_logs.csv")
	c.Header("Content-Type", "text/csv")

	w := csv.NewWriter(c.Writer)
	defer w.Flush()

	_ = w.Write([]string{"CreatedAt", "UserEmail", "Action", "Entity", "EntityID"})
	for _, log := range logs {
		userEmail := ""
		if log.User != nil {
			userEmail = log.User.Email
		}
		entityID := ""
		if log.EntityID != nil {
			entityID = log.EntityID.String()
		}
		_ = w.Write([]string{
			log.CreatedAt.Format(time.RFC3339),
			userEmail,
			log.Action,
			log.Entity,
			entityID,
		})
	}
}

func (h *AuditHandler) ExportPDF(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "300"))
	from, to, userID, action := parseAuditFilters(c)

	logs, err := h.service.ListFiltered(limit, from, to, userID, action)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=audit_logs.pdf")
	c.Header("Content-Type", "application/pdf")

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(0, 10, "Audit Logs")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(50, 7, "Created At", "1", 0, "L", false, 0, "")
	pdf.CellFormat(70, 7, "User Email", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 7, "Action", "1", 0, "L", false, 0, "")
	pdf.CellFormat(45, 7, "Entity", "1", 0, "L", false, 0, "")
	pdf.CellFormat(65, 7, "Entity ID", "1", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	for _, log := range logs {
		userEmail := ""
		if log.User != nil {
			userEmail = log.User.Email
		}
		entityID := ""
		if log.EntityID != nil {
			entityID = log.EntityID.String()
		}

		pdf.CellFormat(50, 6, log.CreatedAt.Format("2006-01-02 15:04"), "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 6, userEmail, "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 6, log.Action, "1", 0, "L", false, 0, "")
		pdf.CellFormat(45, 6, log.Entity, "1", 0, "L", false, 0, "")
		pdf.CellFormat(65, 6, entityID, "1", 1, "L", false, 0, "")
	}

	_ = pdf.Output(c.Writer)
}

func parseAuditFilters(c *gin.Context) (*time.Time, *time.Time, *string, string) {
	fromStr := c.Query("from")
	toStr := c.Query("to")
	userIDStr := c.Query("user_id")
	action := c.Query("action")

	var from, to *time.Time
	if fromStr != "" {
		t, err := time.Parse("2006-01-02", fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse("2006-01-02", toStr)
		if err == nil {
			to = &t
		}
	}

	var userID *string
	if userIDStr != "" {
		userID = &userIDStr
	}
	return from, to, userID, action
}

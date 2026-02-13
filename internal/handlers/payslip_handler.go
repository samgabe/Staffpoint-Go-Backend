package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"

	"go-backend/internal/authz"
	"go-backend/internal/services"
)

type PayslipHandler struct {
	service services.PayslipService
}

func NewPayslipHandler(service services.PayslipService) *PayslipHandler {
	return &PayslipHandler{service: service}
}

func (h *PayslipHandler) Generate(c *gin.Context) {
	generatedBy, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}

	var req struct {
		EmployeeID string  `json:"employee_id" binding:"required,uuid"`
		Month      int     `json:"month" binding:"required"`
		Year       int     `json:"year" binding:"required"`
		BasicPay   float64 `json:"basic_pay" binding:"required"`
		Allowances float64 `json:"allowances"`
		Deductions float64 `json:"deductions"`
		Currency   string  `json:"currency"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	employeeID, _ := uuid.Parse(req.EmployeeID)
	payslip, err := h.service.Generate(employeeID, req.Month, req.Year, req.BasicPay, req.Allowances, req.Deductions, req.Currency, generatedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payslip)
}

func (h *PayslipHandler) Mine(c *gin.Context) {
	userID, err := uuid.Parse(c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	payslips, err := h.service.ListMine(userID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payslips)
}

func (h *PayslipHandler) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "200"))
	var employeeID *uuid.UUID
	if v := c.Query("employee_id"); v != "" {
		if parsed, err := uuid.Parse(v); err == nil {
			employeeID = &parsed
		}
	}
	var month *int
	if v := c.Query("month"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			month = &parsed
		}
	}
	var year *int
	if v := c.Query("year"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil {
			year = &parsed
		}
	}
	payslips, err := h.service.ListAll(limit, employeeID, month, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, payslips)
}

func (h *PayslipHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payslip id"})
		return
	}

	payslip, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payslip not found"})
		return
	}

	role := c.GetString("role")
	userID, _ := uuid.Parse(c.GetString("user_id"))
	if !authz.HasPermission(authz.PermissionsForRole(role), authz.PermManagePayslips) && payslip.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":                "forbidden",
			"code":                 "FORBIDDEN",
			"required_roles":       []string{},
			"required_permissions": []string{authz.PermViewOwnPayslips},
		})
		return
	}

	c.JSON(http.StatusOK, payslip)
}

func (h *PayslipHandler) DownloadPDF(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payslip id"})
		return
	}

	payslip, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "payslip not found"})
		return
	}

	role := c.GetString("role")
	userID, _ := uuid.Parse(c.GetString("user_id"))
	canManage := authz.HasPermission(authz.PermissionsForRole(role), authz.PermManagePayslips)
	if !canManage && payslip.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error":                "forbidden",
			"code":                 "FORBIDDEN",
			"required_roles":       []string{},
			"required_permissions": []string{authz.PermViewOwnPayslips},
		})
		return
	}

	c.Header("Content-Disposition", "attachment; filename=payslip_"+id.String()+".pdf")
	c.Header("Content-Type", "application/pdf")

	employeeName := (payslip.Employee.FirstName + " " + payslip.Employee.LastName)
	if employeeName == " " {
		employeeName = "Employee"
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(15, 15, 15)

	// Header band
	pdf.SetFillColor(17, 44, 99)
	pdf.Rect(15, 15, 180, 24, "F")
	pdf.SetTextColor(255, 255, 255)
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(20, 29, "StaffPoint Payslip")
	pdf.SetFont("Arial", "", 10)
	pdf.Text(20, 35, "Confidential Payroll Document")

	// Meta section
	pdf.SetTextColor(15, 23, 42)
	pdf.SetFillColor(248, 250, 252)
	pdf.Rect(15, 43, 180, 30, "F")
	pdf.SetFont("Arial", "B", 10)
	pdf.Text(20, 51, "Employee")
	pdf.Text(100, 51, "Pay Period")
	pdf.SetFont("Arial", "", 11)
	pdf.Text(20, 58, employeeName)
	pdf.Text(100, 58, monthYearLabel(payslip.Month, payslip.Year))
	pdf.SetFont("Arial", "", 9)
	pdf.Text(20, 66, "Generated: "+time.Now().Format("2006-01-02 15:04"))
	pdf.Text(100, 66, "Payslip ID: "+id.String())

	currency := payslip.Currency
	pdf.SetY(80)
	pdf.SetFont("Arial", "B", 11)
	pdf.SetFillColor(226, 232, 240)
	pdf.CellFormat(120, 8, "Component", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 8, "Amount ("+currency+")", "1", 1, "R", true, 0, "")

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(120, 8, "Basic Pay", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 8, formatCurrency(payslip.BasicPay), "1", 1, "R", false, 0, "")
	pdf.CellFormat(120, 8, "Allowances", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 8, formatCurrency(payslip.Allowances), "1", 1, "R", false, 0, "")
	pdf.CellFormat(120, 8, "Deductions", "1", 0, "L", false, 0, "")
	pdf.CellFormat(60, 8, formatCurrency(payslip.Deductions), "1", 1, "R", false, 0, "")

	// Net pay emphasis
	pdf.SetFont("Arial", "B", 12)
	pdf.SetFillColor(220, 252, 231)
	pdf.CellFormat(120, 10, "Net Pay", "1", 0, "L", true, 0, "")
	pdf.CellFormat(60, 10, formatCurrency(payslip.NetPay), "1", 1, "R", true, 0, "")

	// Footer note
	pdf.Ln(8)
	pdf.SetTextColor(71, 85, 105)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(0, 5, "This payslip is system-generated by StaffPoint and intended for payroll reference only.", "", "L", false)

	_ = pdf.Output(c.Writer)
}

func formatCurrency(value float64) string {
	return strconv.FormatFloat(value, 'f', 2, 64)
}

func monthYearLabel(month, year int) string {
	labels := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	if month >= 1 && month <= 12 {
		return labels[month] + " " + strconv.Itoa(year)
	}
	return strconv.Itoa(month) + "/" + strconv.Itoa(year)
}

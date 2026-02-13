package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-backend/internal/services"
)

type EmployeeHandler struct {
	service *services.EmployeeService
}

func NewEmployeeHandler(service *services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

// POST /employees
func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req struct {
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name" binding:"required"`
		Email        string `json:"email" binding:"required,email"`
		Password     string `json:"password" binding:"required,min=8"`
		Role         string `json:"role" binding:"required,oneof=employee manager"`
		DepartmentID string `json:"department_id" binding:"required,uuid"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deptID, _ := uuid.Parse(req.DepartmentID)
	adminID, _ := uuid.Parse(c.GetString("user_id")) // from JWT

	employee, err := h.service.CreateEmployee(
		req.FirstName,
		req.LastName,
		req.Email,
		req.Role,
		req.Password,
		deptID,
		adminID,
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"employee_id": employee.ID,
		"user_id":     employee.UserID,
		"email":       req.Email,
		"first_name":  employee.FirstName,
		"last_name":   employee.LastName,
	})
}

// List Employee
func (h *EmployeeHandler) ListEmployees(c *gin.Context) {
	employees, err := h.service.ListEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, employees)
}

func (h *EmployeeHandler) CountEmployees(c *gin.Context) {
	count, err := h.service.CountEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"count": count})
}

// PUT /employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	employeeID, _ := uuid.Parse(c.Param("id"))

	var req struct {
		FirstName    string `json:"first_name" binding:"required"`
		LastName     string `json:"last_name" binding:"required"`
		DepartmentID string `json:"department_id" binding:"required,uuid"`
		Role         string `json:"role" binding:"required,oneof=employee manager"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deptID, _ := uuid.Parse(req.DepartmentID)
	adminID, _ := uuid.Parse(c.GetString("user_id"))

	employee, err := h.service.UpdateEmployee(employeeID, req.FirstName, req.LastName, deptID, req.Role, adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"employee_id": employee.ID,
		"first_name":  employee.FirstName,
		"last_name":   employee.LastName,
		"role":        req.Role,
	})
}

// DELETE /employees/:id
func (h *EmployeeHandler) DeactivateEmployee(c *gin.Context) {
	employeeID, _ := uuid.Parse(c.Param("id"))
	adminID, _ := uuid.Parse(c.GetString("user_id"))

	if err := h.service.DeactivateEmployee(employeeID, adminID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Employee deactivated successfully"})
}

package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

type PayslipRepository interface {
	Create(payslip *models.Payslip) error
	Update(payslip *models.Payslip) error
	FindByID(id uuid.UUID) (*models.Payslip, error)
	FindByEmployeePeriod(employeeID uuid.UUID, month, year int) (*models.Payslip, error)
	ListByUser(userID uuid.UUID, limit int) ([]models.Payslip, error)
	ListAll(limit int, employeeID *uuid.UUID, month, year *int) ([]models.Payslip, error)
}

type payslipRepository struct {
	db *gorm.DB
}

func NewPayslipRepository(db *gorm.DB) PayslipRepository {
	return &payslipRepository{db: db}
}

func (r *payslipRepository) Create(payslip *models.Payslip) error {
	return r.db.Create(payslip).Error
}

func (r *payslipRepository) Update(payslip *models.Payslip) error {
	return r.db.Save(payslip).Error
}

func (r *payslipRepository) FindByID(id uuid.UUID) (*models.Payslip, error) {
	var payslip models.Payslip
	if err := r.db.Preload("User").Preload("Employee").First(&payslip, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (r *payslipRepository) FindByEmployeePeriod(employeeID uuid.UUID, month, year int) (*models.Payslip, error) {
	var payslip models.Payslip
	if err := r.db.Where("employee_id = ? AND month = ? AND year = ?", employeeID, month, year).First(&payslip).Error; err != nil {
		return nil, err
	}
	return &payslip, nil
}

func (r *payslipRepository) ListByUser(userID uuid.UUID, limit int) ([]models.Payslip, error) {
	if limit <= 0 {
		limit = 50
	}
	var payslips []models.Payslip
	err := r.db.Preload("Employee").Where("user_id = ?", userID).Order("year DESC, month DESC, created_at DESC").Limit(limit).Find(&payslips).Error
	return payslips, err
}

func (r *payslipRepository) ListAll(limit int, employeeID *uuid.UUID, month, year *int) ([]models.Payslip, error) {
	if limit <= 0 {
		limit = 100
	}
	db := r.db.Preload("Employee").Preload("User").Order("year DESC, month DESC, created_at DESC").Limit(limit)
	if employeeID != nil {
		db = db.Where("employee_id = ?", *employeeID)
	}
	if month != nil {
		db = db.Where("month = ?", *month)
	}
	if year != nil {
		db = db.Where("year = ?", *year)
	}
	var payslips []models.Payslip
	err := db.Find(&payslips).Error
	return payslips, err
}

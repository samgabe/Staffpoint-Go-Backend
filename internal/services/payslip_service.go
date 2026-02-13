package services

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type PayslipService interface {
	Generate(employeeID uuid.UUID, month, year int, basicPay, allowances, deductions float64, currency string, generatedBy uuid.UUID) (*models.Payslip, error)
	ListMine(userID uuid.UUID, limit int) ([]models.Payslip, error)
	ListAll(limit int, employeeID *uuid.UUID, month, year *int) ([]models.Payslip, error)
	GetByID(id uuid.UUID) (*models.Payslip, error)
}

type payslipService struct {
	payslipRepo  repositories.PayslipRepository
	employeeRepo repositories.EmployeeRepository
	auditSvc     AuditService
}

func NewPayslipService(payslipRepo repositories.PayslipRepository, employeeRepo repositories.EmployeeRepository, auditSvc AuditService) PayslipService {
	return &payslipService{payslipRepo: payslipRepo, employeeRepo: employeeRepo, auditSvc: auditSvc}
}

func (s *payslipService) Generate(employeeID uuid.UUID, month, year int, basicPay, allowances, deductions float64, currency string, generatedBy uuid.UUID) (*models.Payslip, error) {
	if month < 1 || month > 12 {
		return nil, errors.New("month must be between 1 and 12")
	}
	if year < 2000 || year > 2100 {
		return nil, errors.New("invalid year")
	}
	if currency == "" {
		currency = "USD"
	}

	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return nil, err
	}

	net := basicPay + allowances - deductions
	if net < 0 {
		net = 0
	}

	existing, err := s.payslipRepo.FindByEmployeePeriod(employeeID, month, year)
	if err == nil {
		existing.BasicPay = basicPay
		existing.Allowances = allowances
		existing.Deductions = deductions
		existing.NetPay = net
		existing.Currency = currency
		existing.GeneratedBy = &generatedBy
		if saveErr := s.payslipRepo.Update(existing); saveErr != nil {
			return nil, saveErr
		}
		s.auditSvc.Log(generatedBy, "PAYSLIP_UPDATED", "payslip", &existing.ID, map[string]interface{}{"month": month, "year": year})
		return existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	payslip := &models.Payslip{
		UserID:      employee.UserID,
		EmployeeID:  employeeID,
		Month:       month,
		Year:        year,
		BasicPay:    basicPay,
		Allowances:  allowances,
		Deductions:  deductions,
		NetPay:      net,
		Currency:    currency,
		GeneratedBy: &generatedBy,
	}

	if err := s.payslipRepo.Create(payslip); err != nil {
		return nil, err
	}
	s.auditSvc.Log(generatedBy, "PAYSLIP_CREATED", "payslip", &payslip.ID, map[string]interface{}{"month": month, "year": year})
	return payslip, nil
}

func (s *payslipService) ListMine(userID uuid.UUID, limit int) ([]models.Payslip, error) {
	return s.payslipRepo.ListByUser(userID, limit)
}

func (s *payslipService) ListAll(limit int, employeeID *uuid.UUID, month, year *int) ([]models.Payslip, error) {
	return s.payslipRepo.ListAll(limit, employeeID, month, year)
}

func (s *payslipService) GetByID(id uuid.UUID) (*models.Payslip, error) {
	return s.payslipRepo.FindByID(id)
}

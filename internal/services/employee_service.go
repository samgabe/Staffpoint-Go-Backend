package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type EmployeeService struct {
	userRepo      repositories.UserRepository
	employeeRepo  repositories.EmployeeRepository
	auditSvc      AuditService
}

func NewEmployeeService(
	userRepo repositories.UserRepository,
	employeeRepo repositories.EmployeeRepository,
	auditSvc AuditService,
) *EmployeeService {
	return &EmployeeService{
		userRepo: userRepo,
		employeeRepo: employeeRepo,
		auditSvc: auditSvc,
	}
}

// CreateEmployee creates a new user + employee record (Admin only)
func (s *EmployeeService) CreateEmployee(
	firstName, lastName, email, role, password string,
	departmentID uuid.UUID,
	adminID uuid.UUID, // for audit logging
) (*models.Employee, error) {

	// 1. Check if user already exists
	existingUser, _ := s.userRepo.FindByEmail(email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// 2. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Create user record
	user := &models.User{
		Email:       email,
		PasswordHash: string(hashedPassword), 
		Role:        role,
		IsActive:    true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 4. Create employee record
	employee := &models.Employee{
		UserID:       user.ID,
		FirstName:    firstName,     
		LastName:     lastName,       
		DepartmentID: &departmentID,  
		Status:       "active",       
		HireDate:     time.Now().UTC(), 
	}

	if err := s.employeeRepo.Create(employee); err != nil {
		return nil, err
	}

	// 5. Audit log
	s.auditSvc.Log(adminID, "EMPLOYEE_CREATED", "employee", &employee.ID, map[string]interface{}{
		"email": email,
		"role":  role,
	})

	return employee, nil
}

// ListEmployees retrieves all employees (Admin/Manager only)
func (s *EmployeeService) ListEmployees() ([]*models.Employee, error) {
	employees, err := s.employeeRepo.List()
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (s *EmployeeService) CountEmployees() (int64, error) {
	return s.employeeRepo.Count()
}


// UpdateEmployee updates employee details (Admin only)
func (s *EmployeeService) UpdateEmployee(
	employeeID uuid.UUID,
	firstName, lastName string,
	departmentID uuid.UUID,
	role string,
	adminID uuid.UUID, // for audit
) (*models.Employee, error) {

	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return nil, err
	}

	employee.FirstName = firstName
	employee.LastName = lastName
	employee.DepartmentID = &departmentID

	// Update user role
	user, err := s.userRepo.FindByID(employee.UserID)
	if err != nil {
		return nil, err
	}
	user.Role = role

	// Save updates
	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	// Audit log
	s.auditSvc.Log(adminID, "EMPLOYEE_UPDATED", "employee", &employee.ID, map[string]interface{}{
		"first_name": firstName,
		"last_name":  lastName,
		"role":       role,
	})

	return employee, nil
}

// DeactivateEmployee sets employee as inactive (Admin only)
func (s *EmployeeService) DeactivateEmployee(
	employeeID uuid.UUID,
	adminID uuid.UUID,
) error {
	employee, err := s.employeeRepo.FindByID(employeeID)
	if err != nil {
		return err
	}

	employee.Status = "inactive"

	if err := s.employeeRepo.Update(employee); err != nil {
		return err
	}

	// Audit log
	s.auditSvc.Log(adminID, "EMPLOYEE_DEACTIVATED", "employee", &employee.ID, nil)

	return nil
}


package services

import (
	"golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type ProfileService struct {
	userRepo     repositories.UserRepository
	employeeRepo repositories.EmployeeRepository
	auditSvc     AuditService
}

// NewProfileService
func NewProfileService(
	userRepo repositories.UserRepository,
	employeeRepo repositories.EmployeeRepository,
	auditSvc AuditService,
) *ProfileService {
	return &ProfileService{
		userRepo: userRepo,
		employeeRepo: employeeRepo,
		auditSvc: auditSvc,
	}
}

// GetProfile returns user + employee info
func (s *ProfileService) GetProfile(userID uuid.UUID) (*models.Employee, error) {
	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

// UpdateProfile allows self-service updates (name, password)
func (s *ProfileService) UpdateProfile(
	userID uuid.UUID,
	firstName, lastName, newPassword string,
) (*models.Employee, error) {

	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, err
	}

	// Update names
	employee.FirstName = firstName
	employee.LastName = lastName

	// Update password if provided
	if newPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashed)
		if err := s.userRepo.Update(user); err != nil {
			return nil, err
		}
	}

	if err := s.employeeRepo.Update(employee); err != nil {
		return nil, err
	}

	// Audit log
	s.auditSvc.Log(userID, "PROFILE_UPDATED", "employee", &employee.ID, nil)

	return employee, nil
}

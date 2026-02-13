package repositories

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

// Interface
type EmployeeRepository interface {
	Create(employee *models.Employee) error
	Update(employee *models.Employee) error
	FindByID(id uuid.UUID) (*models.Employee, error)
	FindByUserID(userID uuid.UUID) (*models.Employee, error)
	List() ([]*models.Employee, error)
	Count() (int64, error)
}

// Implementation
type employeeRepository struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db}
}

func (r *employeeRepository) Create(employee *models.Employee) error {
	return r.db.Create(employee).Error
}

func (r *employeeRepository) Update(employee *models.Employee) error {
	return r.db.Save(employee).Error
}

func (r *employeeRepository) FindByID(id uuid.UUID) (*models.Employee, error) {
	var emp models.Employee
	if err := r.db.Preload("User").Preload("Department").First(&emp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) FindByUserID(userID uuid.UUID) (*models.Employee, error) {
	var emp models.Employee
	if err := r.db.Preload("User").Preload("Department").First(&emp, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) List() ([]*models.Employee, error) {
	var employees []*models.Employee
	if err := r.db.Preload("User").Preload("Department").Find(&employees).Error; err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *employeeRepository) Count() (int64, error) {
	var count int64
	if err := r.db.Model(&models.Employee{}).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}


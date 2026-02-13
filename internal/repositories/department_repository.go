package repositories

import (
	"go-backend/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepository interface {
	Create(dept *models.Department) error
	List() ([]models.Department, error)
	FindByID(id string) (*models.Department, error)
}

type departmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(dept *models.Department) error {
	return r.db.Create(dept).Error
}

func (r *departmentRepository) List() ([]models.Department, error) {
	var departments []models.Department
	err := r.db.Find(&departments).Error
	return departments, err
}

func (r *departmentRepository) FindByID(id string) (*models.Department, error) {
	var dept models.Department
	if err := r.db.First(&dept, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &dept, nil
}

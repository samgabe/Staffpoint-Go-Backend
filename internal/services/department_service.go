package services

import (
	"errors"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type DepartmentService interface {
	Create(name string) (*models.Department, error)
	List() ([]models.Department, error)
}

type departmentService struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentService(repo repositories.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) Create(name string) (*models.Department, error) {
	if name == "" {
		return nil, errors.New("department name is required")
	}

	dept := &models.Department{
		Name: name,
	}

	if err := s.repo.Create(dept); err != nil {
		return nil, err
	}

	return dept, nil
}

func (s *departmentService) List() ([]models.Department, error) {
	return s.repo.List()
}

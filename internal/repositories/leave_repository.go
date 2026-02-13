package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

type LeaveRepository interface {
	Create(req *models.LeaveRequest) error
	Update(req *models.LeaveRequest) error
	FindByID(id uuid.UUID) (*models.LeaveRequest, error)
	ListByUser(userID uuid.UUID, limit int) ([]models.LeaveRequest, error)
	ListAll(status string, limit int) ([]models.LeaveRequest, error)
	CountPending() (int64, error)
}

type leaveRepository struct {
	db *gorm.DB
}

func NewLeaveRepository(db *gorm.DB) LeaveRepository {
	return &leaveRepository{db: db}
}

func (r *leaveRepository) Create(req *models.LeaveRequest) error {
	return r.db.Create(req).Error
}

func (r *leaveRepository) Update(req *models.LeaveRequest) error {
	return r.db.Save(req).Error
}

func (r *leaveRepository) FindByID(id uuid.UUID) (*models.LeaveRequest, error) {
	var leave models.LeaveRequest
	if err := r.db.Preload("User").Preload("Employee").First(&leave, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &leave, nil
}

func (r *leaveRepository) ListByUser(userID uuid.UUID, limit int) ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest
	if limit <= 0 {
		limit = 50
	}

	err := r.db.
		Preload("User").
		Preload("Employee").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&leaves).Error
	return leaves, err
}

func (r *leaveRepository) ListAll(status string, limit int) ([]models.LeaveRequest, error) {
	var leaves []models.LeaveRequest
	if limit <= 0 {
		limit = 50
	}

	db := r.db.Preload("User").Preload("Employee").Order("created_at DESC").Limit(limit)
	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.Find(&leaves).Error
	return leaves, err
}

func (r *leaveRepository) CountPending() (int64, error) {
	var count int64
	err := r.db.Model(&models.LeaveRequest{}).Where("status = ?", "pending").Count(&count).Error
	return count, err
}

func normalizeDate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}

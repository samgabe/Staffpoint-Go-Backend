package repositories

import (
	"time"

	"gorm.io/gorm"

	"go-backend/internal/models"
)

type AuditRepository interface {
	Create(log *models.AuditLog) error
	List(limit int) ([]models.AuditLog, error)
	ListFiltered(limit int, from, to *time.Time, userID *string, action string) ([]models.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) AuditRepository {
	return &auditRepository{db}
}

func (r *auditRepository) Create(log *models.AuditLog) error {
	return r.db.Create(log).Error
}

func (r *auditRepository) List(limit int) ([]models.AuditLog, error) {
	var logs []models.AuditLog

	err := r.db.
		Order("created_at DESC").
		Limit(limit).
		Preload("User").
		Find(&logs).Error

	return logs, err
}
func (r *auditRepository) ListFiltered(limit int, from, to *time.Time, userID *string, action string) ([]models.AuditLog, error) {
	var logs []models.AuditLog
	db := r.db.Order("created_at DESC").Limit(limit).Preload("User")

	if from != nil {
		db = db.Where("created_at >= ?", *from)
	}
	if to != nil {
		endOfDay := to.UTC().Add(24*time.Hour - time.Nanosecond)
		db = db.Where("created_at <= ?", endOfDay)
	}
	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}
	if action != "" {
		db = db.Where("action = ?", action)
	}

	err := db.Find(&logs).Error
	return logs, err
}

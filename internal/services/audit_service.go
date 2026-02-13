package services

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

// AuditService defines the methods for audit logging
type AuditService interface {
	Log(
		userID uuid.UUID,
		action string,
		entity string,
		entityID *uuid.UUID, // pointer matches model
		metadata map[string]interface{},
	)
	List(limit int) ([]models.AuditLog, error)
	ListFiltered(limit int, from, to *time.Time, userID *string, action string) ([]models.AuditLog, error)
}

type auditService struct {
	repo repositories.AuditRepository
}

func NewAuditService(repo repositories.AuditRepository) AuditService {
	return &auditService{repo}
}

// Log creates an audit entry. Fire-and-forget; errors are ignored
func (s *auditService) Log(
	userID uuid.UUID,
	action string,
	entity string,
	entityID *uuid.UUID,
	metadata map[string]interface{},
) {
	var meta datatypes.JSON = []byte("{}") // default empty JSON

	if metadata != nil {
		b, err := json.Marshal(metadata) // standard JSON marshal
		if err == nil {
			meta = datatypes.JSON(b)
		}
	}

	log := &models.AuditLog{
		UserID:   &userID,
		Action:   action,
		Entity:   entity,
		EntityID: entityID,
		Metadata: meta,
	}

	// Fire-and-forget: do not break business logic if audit fails
	_ = s.repo.Create(log)
}
func (s *auditService) List(limit int) ([]models.AuditLog, error) {
	return s.repo.List(limit)
}
func (s *auditService) ListFiltered(limit int, from, to *time.Time, userID *string, action string) ([]models.AuditLog, error) {
	return s.repo.ListFiltered(limit, from, to, userID, action)
}


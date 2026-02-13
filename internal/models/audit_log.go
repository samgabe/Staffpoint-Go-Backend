package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditLog struct {
	BaseModel

	UserID   *uuid.UUID
	Action   string `gorm:"type:varchar(100);not null"`
	Entity   string `gorm:"type:varchar(100);not null"`
	EntityID *uuid.UUID
	Metadata datatypes.JSON

	User *User
}

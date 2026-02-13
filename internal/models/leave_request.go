package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaveRequest struct {
	BaseModel

	UserID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	EmployeeID uuid.UUID  `gorm:"type:uuid;not null;index"`
	StartDate  time.Time  `gorm:"not null"`
	EndDate    time.Time  `gorm:"not null"`
	Reason     string     `gorm:"type:text;not null"`
	Status     string     `gorm:"type:varchar(30);not null;default:'pending';index"`
	ReviewedBy *uuid.UUID `gorm:"type:uuid"`
	ReviewedAt *time.Time

	User     User
	Employee Employee
}

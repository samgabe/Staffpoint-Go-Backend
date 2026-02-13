package models

import (
	"time"

	"github.com/google/uuid"
	
	)
type Employee struct {
	BaseModel

	UserID       uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	DepartmentID *uuid.UUID
	Status       string `gorm:"type:varchar(50);not null"`
	HireDate     time.Time
	FirstName    string `gorm:"type:varchar(100);not null"` // add
	LastName     string `gorm:"type:varchar(100);not null"` // add

	User        User
	Department  *Department
	Attendances []Attendance
}

	
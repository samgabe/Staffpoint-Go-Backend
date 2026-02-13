package models

import (
	"time"

	"github.com/google/uuid"
)

type Attendance struct {
	BaseModel

	EmployeeID uuid.UUID `gorm:"type:uuid;not null;index"`
	WorkDate   time.Time `gorm:"type:date;not null"`
	ClockIn    *time.Time
	ClockOut   *time.Time

	Employee Employee
}

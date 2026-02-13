package models

type User struct {
	BaseModel

	Email        string `gorm:"uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`
	Role         string `gorm:"type:varchar(50);not null"`
	IsActive     bool   `gorm:"default:true"`

	Employee  *Employee
	AuditLogs []AuditLog
}


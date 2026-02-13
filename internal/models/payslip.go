package models

import "github.com/google/uuid"

type Payslip struct {
	BaseModel

	UserID      uuid.UUID  `gorm:"type:uuid;not null;index"`
	EmployeeID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_payslip_employee_period,unique"`
	Month       int        `gorm:"not null;index:idx_payslip_employee_period,unique"`
	Year        int        `gorm:"not null;index:idx_payslip_employee_period,unique"`
	BasicPay    float64    `gorm:"not null"`
	Allowances  float64    `gorm:"not null"`
	Deductions  float64    `gorm:"not null"`
	NetPay      float64    `gorm:"not null"`
	Currency    string     `gorm:"type:varchar(10);not null;default:'USD'"`
	GeneratedBy *uuid.UUID `gorm:"type:uuid"`

	User     User
	Employee Employee
}

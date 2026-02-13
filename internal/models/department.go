package models

type Department struct {
	BaseModel

	Name string `gorm:"uniqueIndex;not null"`

	Employees []Employee
}

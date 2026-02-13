package repositories

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

type AttendanceRepository interface {
	FindByEmployeeAndDate(employeeID uuid.UUID, date time.Time) (*models.Attendance, error)
	FindByDate(date time.Time) ([]models.Attendance, error)
	FindBetweenDates(from, to time.Time) ([]models.Attendance, error)
	Create(attendance *models.Attendance) error
	Update(attendance *models.Attendance) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db}
}

func (r *attendanceRepository) FindByEmployeeAndDate(employeeID uuid.UUID, date time.Time) (*models.Attendance, error) {
	var attendance models.Attendance
	err := r.db.
		Where("employee_id = ? AND work_date = ?", employeeID, date).
		First(&attendance).Error

	return &attendance, err
}

func (r *attendanceRepository) FindByDate(date time.Time) ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.db.
		Preload("Employee").
		Where("work_date = ?", date).
		Find(&records).Error
	return records, err
}

func (r *attendanceRepository) FindBetweenDates(from, to time.Time) ([]models.Attendance, error) {
	var records []models.Attendance
	err := r.db.
		Preload("Employee").
		Where("work_date BETWEEN ? AND ?", from, to).
		Find(&records).Error
	return records, err
}

func (r *attendanceRepository) Create(a *models.Attendance) error {
	return r.db.Create(a).Error
}

func (r *attendanceRepository) Update(a *models.Attendance) error {
	return r.db.Save(a).Error
}

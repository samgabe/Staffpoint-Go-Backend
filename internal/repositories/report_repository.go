package repositories

import (
	"time"

	"gorm.io/gorm"
)

type AttendanceReportRow struct {
	Date      time.Time
	Email     string
	ClockIn   *time.Time
	ClockOut  *time.Time
}

type ReportRepository interface {
	AttendanceReport(from, to time.Time) ([]AttendanceReportRow, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) AttendanceReport(from, to time.Time) ([]AttendanceReportRow, error) {
	var rows []AttendanceReportRow

	err := r.db.Raw(`
		SELECT 
			a.work_date AS date,
			u.email,
			a.clock_in,
			a.clock_out
		FROM attendances a
		JOIN employees e ON a.employee_id = e.id
		JOIN users u ON e.user_id = u.id
		WHERE a.work_date BETWEEN ? AND ?
		ORDER BY a.work_date, u.email
	`, from, to).Scan(&rows).Error

	return rows, err
}

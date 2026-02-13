package databases

import (
	"gorm.io/gorm"

	"go-backend/internal/models"
)

func Migrate(db *gorm.DB) error {
	// Auto-migrate tables
	if err := db.AutoMigrate(
		&models.User{},
		&models.Department{},
		&models.Employee{},
		&models.Attendance{},
		&models.AuditLog{},
		&models.LeaveRequest{},
		&models.Payslip{},
	); err != nil {
		return err
	}

	// Attendance indexes
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_attendance_employee_date
		ON attendances (employee_id, work_date)
	`)

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_attendance_work_date
		ON attendances (work_date)
	`)

	// Audit logs index
	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at
		ON audit_logs (created_at)
	`)

	db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_leave_requests_status
		ON leave_requests (status)
	`)

	return nil
}

package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type AttendanceService struct {
	attendanceRepo repositories.AttendanceRepository
	employeeRepo   repositories.EmployeeRepository
	auditSvc       AuditService
}

func NewAttendanceService(
	attendanceRepo repositories.AttendanceRepository,
	employeeRepo repositories.EmployeeRepository,
	auditSvc AuditService,
) *AttendanceService {
	return &AttendanceService{
		attendanceRepo: attendanceRepo,
		employeeRepo:   employeeRepo,
		auditSvc:       auditSvc,
	}
}

func (s *AttendanceService) ClockInByUser(userID uuid.UUID) error {

	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("employee profile not found")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	_, err = s.attendanceRepo.FindByEmployeeAndDate(employee.ID, today)
	if err == nil {
		return errors.New("already clocked in today")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	now := time.Now().UTC()
	attendance := &models.Attendance{
		EmployeeID: employee.ID,
		WorkDate:   today,
		ClockIn:    &now,
	}

	if err := s.attendanceRepo.Create(attendance); err != nil {
		return err
	}

	s.auditSvc.Log(userID, "CLOCK_IN", "attendance", &attendance.ID, nil)

	return nil
}

func (s *AttendanceService) ClockOutByUser(userID uuid.UUID) error {

	employee, err := s.employeeRepo.FindByUserID(userID)
	if err != nil {
		return errors.New("employee profile not found")
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	attendance, err := s.attendanceRepo.FindByEmployeeAndDate(employee.ID, today)
	if err != nil {
		return errors.New("no active attendance record")
	}

	if attendance.ClockIn == nil {
		return errors.New("invalid attendance state")
	}

	if attendance.ClockOut != nil {
		return errors.New("already clocked out")
	}

	now := time.Now().UTC()
	if now.Before(*attendance.ClockIn) {
		return errors.New("clock-out before clock-in")
	}

	attendance.ClockOut = &now

	if err := s.attendanceRepo.Update(attendance); err != nil {
		return err
	}

	s.auditSvc.Log(userID, "CLOCK_OUT", "attendance", &attendance.ID, nil)

	return nil
}

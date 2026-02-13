package services

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type LeaveService interface {
	RequestLeave(userID, employeeID uuid.UUID, startDate, endDate time.Time, reason string) (*models.LeaveRequest, error)
	ListMine(userID uuid.UUID, limit int) ([]models.LeaveRequest, error)
	ListAll(status string, limit int) ([]models.LeaveRequest, error)
	ReviewLeave(leaveID, reviewerID uuid.UUID, status string) (*models.LeaveRequest, error)
	CancelMyLeave(leaveID, userID uuid.UUID) (*models.LeaveRequest, error)
	PendingCount() (int64, error)
}

type leaveService struct {
	repo     repositories.LeaveRepository
	auditSvc AuditService
}

func NewLeaveService(repo repositories.LeaveRepository, auditSvc AuditService) LeaveService {
	return &leaveService{repo: repo, auditSvc: auditSvc}
}

func (s *leaveService) RequestLeave(userID, employeeID uuid.UUID, startDate, endDate time.Time, reason string) (*models.LeaveRequest, error) {
	if strings.TrimSpace(reason) == "" {
		return nil, errors.New("reason is required")
	}
	if endDate.Before(startDate) {
		return nil, errors.New("end_date must be on or after start_date")
	}

	leave := &models.LeaveRequest{
		UserID:     userID,
		EmployeeID: employeeID,
		StartDate:  startDate.UTC(),
		EndDate:    endDate.UTC(),
		Reason:     strings.TrimSpace(reason),
		Status:     "pending",
	}

	if err := s.repo.Create(leave); err != nil {
		return nil, err
	}

	s.auditSvc.Log(userID, "LEAVE_REQUESTED", "leave_request", &leave.ID, map[string]interface{}{
		"start_date": leave.StartDate.Format("2006-01-02"),
		"end_date":   leave.EndDate.Format("2006-01-02"),
	})

	return leave, nil
}

func (s *leaveService) ListMine(userID uuid.UUID, limit int) ([]models.LeaveRequest, error) {
	return s.repo.ListByUser(userID, limit)
}

func (s *leaveService) ListAll(status string, limit int) ([]models.LeaveRequest, error) {
	return s.repo.ListAll(status, limit)
}

func (s *leaveService) ReviewLeave(leaveID, reviewerID uuid.UUID, status string) (*models.LeaveRequest, error) {
	normalized := strings.ToLower(strings.TrimSpace(status))
	if normalized != "approved" && normalized != "rejected" {
		return nil, errors.New("status must be approved or rejected")
	}

	leave, err := s.repo.FindByID(leaveID)
	if err != nil {
		return nil, err
	}
	if leave.UserID == reviewerID {
		return nil, errors.New("you cannot review your own leave request")
	}

	if leave.Status != "pending" {
		return nil, errors.New("leave request has already been reviewed")
	}

	now := time.Now().UTC()
	leave.Status = normalized
	leave.ReviewedBy = &reviewerID
	leave.ReviewedAt = &now

	if err := s.repo.Update(leave); err != nil {
		return nil, err
	}

	s.auditSvc.Log(reviewerID, "LEAVE_"+strings.ToUpper(normalized), "leave_request", &leave.ID, nil)
	return leave, nil
}

func (s *leaveService) CancelMyLeave(leaveID, userID uuid.UUID) (*models.LeaveRequest, error) {
	leave, err := s.repo.FindByID(leaveID)
	if err != nil {
		return nil, err
	}
	if leave.UserID != userID {
		return nil, errors.New("you can only cancel your own leave requests")
	}
	if leave.Status != "pending" {
		return nil, errors.New("only pending leave requests can be cancelled")
	}

	now := time.Now().UTC()
	leave.Status = "cancelled"
	leave.ReviewedBy = &userID
	leave.ReviewedAt = &now

	if err := s.repo.Update(leave); err != nil {
		return nil, err
	}

	s.auditSvc.Log(userID, "LEAVE_CANCELLED", "leave_request", &leave.ID, nil)
	return leave, nil
}

func (s *leaveService) PendingCount() (int64, error) {
	return s.repo.CountPending()
}

package services

import (
	"time"

	"github.com/google/uuid"

	"go-backend/internal/authz"
	"go-backend/internal/models"
	"go-backend/internal/repositories"
)

type NotificationItem struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationPayload struct {
	UnreadCount int64              `json:"unread_count"`
	Items       []NotificationItem `json:"items"`
}

type NotificationService interface {
	GetNotifications(userID uuid.UUID, role string) (NotificationPayload, error)
}

type notificationService struct {
	leaveRepo repositories.LeaveRepository
}

func NewNotificationService(leaveRepo repositories.LeaveRepository) NotificationService {
	return &notificationService{leaveRepo: leaveRepo}
}

func (s *notificationService) GetNotifications(userID uuid.UUID, role string) (NotificationPayload, error) {
	payload := NotificationPayload{
		UnreadCount: 0,
		Items:       make([]NotificationItem, 0),
	}

	isReviewer := authz.HasPermission(authz.PermissionsForRole(role), authz.PermReviewLeaves)
	if isReviewer {
		pendingCount, err := s.leaveRepo.CountPending()
		if err != nil {
			return payload, err
		}
		payload.UnreadCount = pendingCount

		leaves, err := s.leaveRepo.ListAll("pending", 5)
		if err != nil {
			return payload, err
		}

		for _, leave := range leaves {
			fullName := (leave.Employee.FirstName + " " + leave.Employee.LastName)
			payload.Items = append(payload.Items, NotificationItem{
				ID:        leave.ID.String(),
				Title:     "Leave Request Pending",
				Message:   fullName + " requested leave",
				Type:      "leave",
				CreatedAt: leave.CreatedAt,
			})
		}
		return payload, nil
	}

	leaves, err := s.leaveRepo.ListByUser(userID, 5)
	if err != nil {
		return payload, err
	}

	for _, leave := range leaves {
		if leave.Status == "pending" {
			payload.UnreadCount++
		}
		payload.Items = append(payload.Items, NotificationItem{
			ID:        leave.ID.String(),
			Title:     "Leave " + leave.Status,
			Message:   leave.StartDate.Format("2006-01-02") + " to " + leave.EndDate.Format("2006-01-02"),
			Type:      "leave",
			CreatedAt: leave.CreatedAt,
		})
	}

	return payload, nil
}

func leaveDisplayName(employee models.Employee) string {
	if employee.FirstName == "" && employee.LastName == "" {
		return "Employee"
	}
	return employee.FirstName + " " + employee.LastName
}

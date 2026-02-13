package services

import (
	"time"

	"go-backend/internal/repositories"
)

type AnalyticsService struct {
	repo repositories.AnalyticsRepository
}

func NewAnalyticsService(repo repositories.AnalyticsRepository) *AnalyticsService {
	return &AnalyticsService{repo}
}

func (s *AnalyticsService) DailySummary(date time.Time) (map[string]int64, error) {
	return s.repo.DailySummary(date)
}

func (s *AnalyticsService) Trend(from, to time.Time) ([]repositories.TrendPoint, error) {
	return s.repo.AttendanceTrend(from, to)
}

func (s *AnalyticsService) Absentees(date time.Time) ([]string, error) {
	return s.repo.Absentees(date)
}

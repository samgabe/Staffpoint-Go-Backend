package repositories

import (
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	DailySummary(date time.Time) (map[string]int64, error)
	AttendanceTrend(from, to time.Time) ([]TrendPoint, error)
	Absentees(date time.Time) ([]string, error)
}

type TrendPoint struct {
	Date  time.Time
	Count int64
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db}
}

func (r *analyticsRepository) DailySummary(date time.Time) (map[string]int64, error) {
	var result struct {
		Present int64
		Absent  int64
	}

	err := r.db.Raw(`
		SELECT
			COUNT(DISTINCT a.employee_id) AS present,
			(
				SELECT COUNT(*) FROM employees e
				WHERE e.status = 'active'
				AND e.id NOT IN (
					SELECT employee_id FROM attendances WHERE work_date = ?
				)
			) AS absent
		FROM attendances a
		WHERE a.work_date = ?
	`, date, date).Scan(&result).Error

	return map[string]int64{
		"present": result.Present,
		"absent":  result.Absent,
	}, err
}

func (r *analyticsRepository) AttendanceTrend(from, to time.Time) ([]TrendPoint, error) {
	var data []TrendPoint

	err := r.db.Raw(`
		SELECT work_date AS date, COUNT(*) AS count
		FROM attendances
		WHERE work_date BETWEEN ? AND ?
		GROUP BY work_date
		ORDER BY work_date
	`, from, to).Scan(&data).Error

	return data, err
}

func (r *analyticsRepository) Absentees(date time.Time) ([]string, error) {
	var emails []string

	err := r.db.Raw(`
		SELECT u.email
		FROM users u
		JOIN employees e ON u.id = e.user_id
		WHERE e.status = 'active'
		AND e.id NOT IN (
			SELECT employee_id FROM attendances WHERE work_date = ?
		)
	`, date).Scan(&emails).Error

	return emails, err
}

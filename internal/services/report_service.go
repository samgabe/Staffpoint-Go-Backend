package services

import (
	"encoding/csv"
	"io"
	"time"
	
	"github.com/jung-kurt/gofpdf"
	
	"go-backend/internal/repositories"
	
)

type ReportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) *ReportService {
	return &ReportService{repo}
}

func (s *ReportService) ExportCSV(w io.Writer, from, to time.Time) error {
	rows, err := s.repo.AttendanceReport(from, to)
	if err != nil {
		return err
	}

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"Date", "Email", "Clock In", "Clock Out"})

	for _, r := range rows {
		writer.Write([]string{
			r.Date.Format("2006-01-02"),
			r.Email,
			formatTime(r.ClockIn),
			formatTime(r.ClockOut),
		})
	}

	return nil
}

func formatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("15:04:05")
}


func (s *ReportService) ExportPDF(w io.Writer, from, to time.Time) error {
	rows, err := s.repo.AttendanceReport(from, to)
	if err != nil {
		return err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Attendance Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 10)

	for _, r := range rows {
		pdf.Cell(0, 8,
			r.Date.Format("2006-01-02")+" | "+
				r.Email+" | "+
				formatTime(r.ClockIn)+" - "+
				formatTime(r.ClockOut),
		)
		pdf.Ln(6)
	}

	return pdf.Output(w)
}

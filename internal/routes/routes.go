package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go-backend/internal/authz"
	"go-backend/internal/handlers"
	"go-backend/internal/middleware"
	"go-backend/internal/repositories"
	"go-backend/internal/services"
)

func RegisterRoutes(router *gin.Engine, db *gorm.DB, jwtSecret string) {
	api := router.Group("/api")

	// ===== Repositories =====
	userRepo := repositories.NewUserRepository(db)
	employeeRepo := repositories.NewEmployeeRepository(db)
	departmentRepo := repositories.NewDepartmentRepository(db)
	attendanceRepo := repositories.NewAttendanceRepository(db)
	analyticsRepo := repositories.NewAnalyticsRepository(db)
	reportRepo := repositories.NewReportRepository(db)
	auditRepo := repositories.NewAuditRepository(db)
	leaveRepo := repositories.NewLeaveRepository(db)
	payslipRepo := repositories.NewPayslipRepository(db)

	// ===== Services =====
	auditSvc := services.NewAuditService(auditRepo)
	authSvc := services.NewAuthService(userRepo, employeeRepo, jwtSecret)
	employeeSvc := services.NewEmployeeService(userRepo, employeeRepo, auditSvc)
	departmentSvc := services.NewDepartmentService(departmentRepo)
	profileSvc := services.NewProfileService(userRepo, employeeRepo, auditSvc)
	attendanceSvc := services.NewAttendanceService(attendanceRepo, employeeRepo, auditSvc)
	analyticsSvc := services.NewAnalyticsService(analyticsRepo)
	reportSvc := services.NewReportService(reportRepo)
	leaveSvc := services.NewLeaveService(leaveRepo, auditSvc)
	notificationSvc := services.NewNotificationService(leaveRepo)
	payslipSvc := services.NewPayslipService(payslipRepo, employeeRepo, auditSvc)
	// Add other services as needed

	// ===== Handlers =====
	authHandler := handlers.NewAuthHandler(authSvc)
	employeeHandler := handlers.NewEmployeeHandler(employeeSvc)
	departmentHandler := handlers.NewDepartmentHandler(departmentSvc)
	profileHandler := handlers.NewProfileHandler(profileSvc)
	attendanceHandler := handlers.NewAttendanceHandler(attendanceSvc)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsSvc)
	reportHandler := handlers.NewReportHandler(reportSvc)
	auditHandler := handlers.NewAuditHandler(auditSvc)
	leaveHandler := handlers.NewLeaveHandler(leaveSvc)
	notificationHandler := handlers.NewNotificationHandler(notificationSvc)
	payslipHandler := handlers.NewPayslipHandler(payslipSvc)

	// ===== Auth Routes =====
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/refresh", authHandler.Refresh)

	// ===== Protected Routes =====
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.Use(middleware.AuditAuthorizationFailures(auditSvc))

	// Employees
	employees := protected.Group("/employees")
	employees.Use(middleware.RequirePermissions(authz.PermManageEmployees))
	employees.GET("/count", employeeHandler.CountEmployees)
	employees.GET("/", employeeHandler.ListEmployees)
	employees.POST("/", employeeHandler.CreateEmployee)
	employees.PUT("/:id", employeeHandler.UpdateEmployee)
	employees.DELETE("/:id", employeeHandler.DeactivateEmployee)

	// Departments
	departments := protected.Group("/departments")
	departments.POST("/", middleware.RequirePermissions(authz.PermManageDepartments), departmentHandler.Create)
	departments.GET("/", middleware.RequirePermissions(authz.PermViewDepartments), departmentHandler.List)

	// Profile
	profile := protected.Group("/profile")
	profile.Use(middleware.RequirePermissions(authz.PermViewProfile))
	profile.GET("/", profileHandler.GetProfile)
	profile.Use(middleware.RequirePermissions(authz.PermUpdateProfile))
	profile.PUT("/", profileHandler.UpdateProfile)

	// Attendance
	attendance := protected.Group("/attendance")
	attendance.Use(middleware.RequirePermissions(authz.PermClockAttendance))
	attendance.POST("/clock-in", attendanceHandler.ClockIn)
	attendance.POST("/clock-out", attendanceHandler.ClockOut)

	// Analytics
	analytics := protected.Group("/analytics")
	analytics.Use(middleware.RequirePermissions(authz.PermViewAnalytics))
	analytics.GET("/daily-summary", analyticsHandler.DailySummary)
	analytics.GET("/attendance-trend", analyticsHandler.Trend)
	analytics.GET("/absentees", analyticsHandler.Absentees)

	// Reports
	reports := protected.Group("/reports")
	reports.Use(middleware.RequirePermissions(authz.PermExportReports))
	reports.GET("/attendance/csv", reportHandler.ExportCSV)
	reports.GET("/attendance/pdf", reportHandler.ExportPDF)

	// Audit logs
	audit := protected.Group("/audit-logs")
	audit.Use(middleware.RequirePermissions(authz.PermViewAuditLogs))
	audit.GET("/", auditHandler.List)
	audit.GET("/csv", auditHandler.ExportCSV)
	audit.GET("/pdf", auditHandler.ExportPDF)

	// Leaves
	leaves := protected.Group("/leaves")
	leaves.Use(middleware.RequirePermissions(authz.PermRequestLeave))
	leaves.POST("/", leaveHandler.Request)
	leaves.Use(middleware.RequirePermissions(authz.PermViewOwnLeaves))
	leaves.GET("/mine", leaveHandler.Mine)
	leaves.PUT("/:id/cancel", leaveHandler.CancelMine)
	leaves.GET("/", middleware.RequirePermissions(authz.PermReviewLeaves), leaveHandler.List)
	leaves.PUT("/:id/review", middleware.RequirePermissions(authz.PermReviewLeaves), leaveHandler.Review)

	// Notifications
	notifications := protected.Group("/notifications")
	notifications.Use(middleware.RequirePermissions(authz.PermViewNotifications))
	notifications.GET("/", notificationHandler.List)

	// Payslips
	payslips := protected.Group("/payslips")
	payslips.Use(middleware.RequirePermissions(authz.PermViewOwnPayslips))
	payslips.GET("/mine", payslipHandler.Mine)
	payslips.GET("/:id", payslipHandler.GetByID)
	payslips.GET("/:id/pdf", payslipHandler.DownloadPDF)
	payslips.GET("/", middleware.RequirePermissions(authz.PermManagePayslips), payslipHandler.List)
	payslips.POST("/", middleware.RequirePermissions(authz.PermManagePayslips), payslipHandler.Generate)
}

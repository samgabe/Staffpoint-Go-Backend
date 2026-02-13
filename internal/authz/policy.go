package authz

const (
	RoleAdmin    = "admin"
	RoleManager  = "manager"
	RoleEmployee = "employee"
)

const (
	PermManageEmployees   = "manage_employees"
	PermManageDepartments = "manage_departments"
	PermViewDepartments   = "view_departments"
	PermViewAnalytics     = "view_analytics"
	PermExportReports     = "export_reports"
	PermViewAuditLogs     = "view_audit_logs"
	PermReviewLeaves      = "review_leaves"
	PermRequestLeave      = "request_leave"
	PermViewOwnLeaves     = "view_own_leaves"
	PermClockAttendance   = "clock_attendance"
	PermViewNotifications = "view_notifications"
	PermViewProfile       = "view_profile"
	PermUpdateProfile     = "update_profile"
	PermManagePayslips    = "manage_payslips"
	PermViewOwnPayslips   = "view_own_payslips"
)

var rolePermissions = map[string][]string{
	RoleAdmin: {
		PermManageEmployees,
		PermManageDepartments,
		PermViewDepartments,
		PermViewAnalytics,
		PermExportReports,
		PermViewAuditLogs,
		PermReviewLeaves,
		PermRequestLeave,
		PermViewOwnLeaves,
		PermClockAttendance,
		PermViewNotifications,
		PermViewProfile,
		PermUpdateProfile,
		PermManagePayslips,
		PermViewOwnPayslips,
	},
	RoleManager: {
		PermManageEmployees,
		PermViewDepartments,
		PermViewAnalytics,
		PermExportReports,
		PermReviewLeaves,
		PermRequestLeave,
		PermViewOwnLeaves,
		PermClockAttendance,
		PermViewNotifications,
		PermViewProfile,
		PermUpdateProfile,
		PermManagePayslips,
		PermViewOwnPayslips,
	},
	RoleEmployee: {
		PermRequestLeave,
		PermViewOwnLeaves,
		PermClockAttendance,
		PermViewNotifications,
		PermViewProfile,
		PermUpdateProfile,
		PermViewOwnPayslips,
	},
}

func PermissionsForRole(role string) []string {
	permissions := rolePermissions[role]
	out := make([]string, len(permissions))
	copy(out, permissions)
	return out
}

func HasRole(role string, allowed ...string) bool {
	for _, candidate := range allowed {
		if role == candidate {
			return true
		}
	}
	return false
}

func HasPermission(permissions []string, required string) bool {
	for _, permission := range permissions {
		if permission == required {
			return true
		}
	}
	return false
}

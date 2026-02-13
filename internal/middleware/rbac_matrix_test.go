package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"go-backend/internal/authz"
	"go-backend/pkg/utils"
)

const testSecret = "rbac-test-secret"

func issueToken(t *testing.T, role string) string {
	t.Helper()
	token, err := utils.GenerateToken(
		"user-1",
		"employee-1",
		role,
		authz.PermissionsForRole(role),
		testSecret,
		time.Hour,
	)
	if err != nil {
		t.Fatalf("failed to issue token: %v", err)
	}
	return token
}

func newRBACRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	api := r.Group("/api")
	api.Use(AuthMiddleware(testSecret))

	api.GET("/employees", RequirePermissions(authz.PermManageEmployees), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.POST("/departments", RequirePermissions(authz.PermManageDepartments), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/analytics", RequirePermissions(authz.PermViewAnalytics), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/reports", RequirePermissions(authz.PermExportReports), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/audit-logs", RequirePermissions(authz.PermViewAuditLogs), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.GET("/leaves", RequirePermissions(authz.PermReviewLeaves), func(c *gin.Context) { c.Status(http.StatusOK) })
	api.POST("/leaves", RequirePermissions(authz.PermRequestLeave), func(c *gin.Context) { c.Status(http.StatusOK) })

	return r
}

func TestRBACPermissionMatrix(t *testing.T) {
	router := newRBACRouter()

	type testCase struct {
		name           string
		role           string
		method         string
		path           string
		expectedStatus int
	}

	cases := []testCase{
		{"admin employees", authz.RoleAdmin, http.MethodGet, "/api/employees", http.StatusOK},
		{"admin departments create", authz.RoleAdmin, http.MethodPost, "/api/departments", http.StatusOK},
		{"admin analytics", authz.RoleAdmin, http.MethodGet, "/api/analytics", http.StatusOK},
		{"admin reports", authz.RoleAdmin, http.MethodGet, "/api/reports", http.StatusOK},
		{"admin audit logs", authz.RoleAdmin, http.MethodGet, "/api/audit-logs", http.StatusOK},
		{"admin review leaves", authz.RoleAdmin, http.MethodGet, "/api/leaves", http.StatusOK},
		{"admin request leaves", authz.RoleAdmin, http.MethodPost, "/api/leaves", http.StatusOK},

		{"manager employees", authz.RoleManager, http.MethodGet, "/api/employees", http.StatusOK},
		{"manager departments create denied", authz.RoleManager, http.MethodPost, "/api/departments", http.StatusForbidden},
		{"manager analytics", authz.RoleManager, http.MethodGet, "/api/analytics", http.StatusOK},
		{"manager reports", authz.RoleManager, http.MethodGet, "/api/reports", http.StatusOK},
		{"manager audit logs denied", authz.RoleManager, http.MethodGet, "/api/audit-logs", http.StatusForbidden},
		{"manager review leaves", authz.RoleManager, http.MethodGet, "/api/leaves", http.StatusOK},
		{"manager request leaves", authz.RoleManager, http.MethodPost, "/api/leaves", http.StatusOK},

		{"employee employees denied", authz.RoleEmployee, http.MethodGet, "/api/employees", http.StatusForbidden},
		{"employee departments create denied", authz.RoleEmployee, http.MethodPost, "/api/departments", http.StatusForbidden},
		{"employee analytics denied", authz.RoleEmployee, http.MethodGet, "/api/analytics", http.StatusForbidden},
		{"employee reports denied", authz.RoleEmployee, http.MethodGet, "/api/reports", http.StatusForbidden},
		{"employee audit logs denied", authz.RoleEmployee, http.MethodGet, "/api/audit-logs", http.StatusForbidden},
		{"employee review leaves denied", authz.RoleEmployee, http.MethodGet, "/api/leaves", http.StatusForbidden},
		{"employee request leaves", authz.RoleEmployee, http.MethodPost, "/api/leaves", http.StatusOK},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			req.Header.Set("Authorization", "Bearer "+issueToken(t, tc.role))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Fatalf("expected %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedStatus == http.StatusForbidden {
				var payload map[string]any
				if err := json.Unmarshal(w.Body.Bytes(), &payload); err != nil {
					t.Fatalf("invalid forbidden payload: %v", err)
				}
				if payload["code"] != "FORBIDDEN" {
					t.Fatalf("expected forbidden code, got %v", payload["code"])
				}
			}
		})
	}
}

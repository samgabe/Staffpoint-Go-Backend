package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-backend/internal/services"
)

func AuditAuthorizationFailures(auditSvc services.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if c.Writer.Status() != http.StatusForbidden {
			return
		}

		userID, err := uuid.Parse(c.GetString("user_id"))
		if err != nil {
			userID = uuid.Nil
		}

		role := c.GetString("role")
		auditSvc.Log(
			userID,
			"AUTHZ_DENIED",
			"route",
			nil,
			map[string]interface{}{
				"path":   c.FullPath(),
				"method": c.Request.Method,
				"role":   role,
			},
		)
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"

	"go-backend/internal/authz"
)

func RequirePermissions(required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		rawPermissions, exists := c.Get("permissions")
		if !exists {
			abortForbidden(c, nil, required)
			return
		}

		permissions, ok := rawPermissions.([]string)
		if !ok {
			abortForbidden(c, nil, required)
			return
		}

		for _, requirement := range required {
			if !authz.HasPermission(permissions, requirement) {
				abortForbidden(c, nil, required)
				return
			}
		}

		c.Next()
	}
}

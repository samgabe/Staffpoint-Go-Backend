package middleware

import (
	"github.com/gin-gonic/gin"

	"go-backend/internal/authz"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")

		if authz.HasRole(userRole, roles...) {
			c.Next()
			return
		}

		abortForbidden(c, roles, nil)
	}
}

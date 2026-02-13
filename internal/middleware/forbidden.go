package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortForbidden(c *gin.Context, requiredRoles []string, requiredPermissions []string) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
		"error":                "forbidden",
		"code":                 "FORBIDDEN",
		"required_roles":       requiredRoles,
		"required_permissions": requiredPermissions,
	})
}

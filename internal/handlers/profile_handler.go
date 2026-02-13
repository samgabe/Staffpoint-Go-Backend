package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"go-backend/internal/authz"
	"go-backend/internal/services"
)

type ProfileHandler struct {
	service *services.ProfileService
}

func NewProfileHandler(service *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{service: service}
}

// GET /profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	profile, err := h.service.GetProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     profile.UserID,
		"first_name":  profile.FirstName,
		"last_name":   profile.LastName,
		"email":       profile.User.Email,
		"role":        profile.User.Role,
		"permissions": authz.PermissionsForRole(profile.User.Role),
		"department":  profile.DepartmentID,
		"status":      profile.Status,
	})
}

// PUT /profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, _ := uuid.Parse(c.GetString("user_id"))

	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Password  string `json:"password,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	profile, err := h.service.UpdateProfile(userID, req.FirstName, req.LastName, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    profile.UserID,
		"first_name": profile.FirstName,
		"last_name":  profile.LastName,
	})
}

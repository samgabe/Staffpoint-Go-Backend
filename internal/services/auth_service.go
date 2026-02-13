package services

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"go-backend/internal/authz"
	"go-backend/internal/repositories"
	"go-backend/pkg/utils"
)

type AuthService struct {
	userRepo     repositories.UserRepository
	employeeRepo repositories.EmployeeRepository
	jwtSecret    string
}

func NewAuthService(
	userRepo repositories.UserRepository,
	employeeRepo repositories.EmployeeRepository,
	secret string,
) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		employeeRepo: employeeRepo,
		jwtSecret:    secret,
	}
}

// Login authenticates a user and returns access & refresh tokens
func (s *AuthService) Login(email, password string) (string, string, error) {
	// 1. Find user by email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// 2. Verify password
	if err := utils.CheckPassword(user.PasswordHash, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// 3. Get employee record (needed for employeeID in JWT)
	employee, err := s.employeeRepo.FindByUserID(user.ID)
	if err != nil {
		return "", "", errors.New("employee record not found")
	}

	// 4. Generate access token
	permissions := authz.PermissionsForRole(user.Role)

	accessToken, err := utils.GenerateToken(
		user.ID.String(),
		employee.ID.String(),
		user.Role,
		permissions,
		s.jwtSecret,
		15*time.Minute,
	)
	if err != nil {
		return "", "", err
	}

	// 5. Generate refresh token
	refreshToken, err := utils.GenerateToken(
		user.ID.String(),
		employee.ID.String(),
		user.Role,
		permissions,
		s.jwtSecret,
		7*24*time.Hour,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Refresh validates a refresh token and rotates access/refresh tokens.
func (s *AuthService) Refresh(refreshToken string) (string, string, error) {
	claims := &utils.JWTClaims{}
	token, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
	)
	if err != nil || !token.Valid {
		return "", "", errors.New("invalid refresh token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return "", "", errors.New("invalid refresh token")
	}

	user, err := s.userRepo.FindByID(userID)
	if err != nil || !user.IsActive {
		return "", "", errors.New("user not found")
	}

	employeeID := claims.EmployeeID
	if employeeID == "" {
		employee, empErr := s.employeeRepo.FindByUserID(user.ID)
		if empErr != nil {
			return "", "", errors.New("employee record not found")
		}
		employeeID = employee.ID.String()
	}

	newAccessToken, err := utils.GenerateToken(
		user.ID.String(),
		employeeID,
		user.Role,
		authz.PermissionsForRole(user.Role),
		s.jwtSecret,
		15*time.Minute,
	)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := utils.GenerateToken(
		user.ID.String(),
		employeeID,
		user.Role,
		authz.PermissionsForRole(user.Role),
		s.jwtSecret,
		7*24*time.Hour,
	)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

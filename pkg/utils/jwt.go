package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID      string   `json:"user_id"`
	Role        string   `json:"role"`
	EmployeeID  string   `json:"employee_id"`
	Permissions []string `json:"permissions,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID, employeeID, role string, permissions []string, secret string, ttl time.Duration) (string, error) {
	claims := JWTClaims{
		UserID:      userID,
		Role:        role,
		EmployeeID:  employeeID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

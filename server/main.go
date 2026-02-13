package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-backend/databases"
	"go-backend/internal/routes"
)

func main() {
	// Environment
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	jwtSecret := os.Getenv("JWT_SECRET")

	if dsn == "" || jwtSecret == "" {
		log.Fatal("Missing required environment variables")
	}

	// Database connection with proper config
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true, // Add this
	})
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Run migrations
	if err := databases.Migrate(db); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Gin router
	router := gin.Default()
	router.Use(corsMiddleware())

	// Register routes
	routes.RegisterRoutes(router, db, jwtSecret)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}

func corsMiddleware() gin.HandlerFunc {
	allowedOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = os.Getenv("CORS_ALLOWED_ORIGIN")
	}
	if allowedOriginsEnv == "" {
		allowedOriginsEnv = "http://localhost:5173,http://127.0.0.1:5173"
	}
	allowedOrigins := splitCSV(allowedOriginsEnv)

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		if isOriginAllowed(origin, allowedOrigins) {
			if contains(allowedOrigins, "*") {
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				c.Writer.Header().Set("Vary", "Origin")
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if strings.EqualFold(c.Request.Method, "OPTIONS") {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}

func contains(values []string, needle string) bool {
	for _, value := range values {
		if value == needle {
			return true
		}
	}
	return false
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	if contains(allowedOrigins, "*") {
		return true
	}
	return contains(allowedOrigins, origin)
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-backend/internal/models"
)

func main() {
	_ = godotenv.Load("../../.env", ".env")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=newman dbname=employee port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	email := getEnv("ADMIN_SEED_EMAIL", "admin@company.com")
	password := getEnv("ADMIN_SEED_PASSWORD", "Admin@123")

	var user models.User
	err = db.Where("email = ?", email).First(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatal(err)
	}

	if err == gorm.ErrRecordNotFound {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			log.Fatal(hashErr)
		}

		user = models.User{
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "admin",
			IsActive:     true,
		}

		if createErr := db.Create(&user).Error; createErr != nil {
			log.Fatal(createErr)
		}
	} else {
		if user.Role != "admin" || !user.IsActive {
			user.Role = "admin"
			user.IsActive = true
			if saveErr := db.Save(&user).Error; saveErr != nil {
				log.Fatal(saveErr)
			}
		}
	}

	var employee models.Employee
	empErr := db.Where("user_id = ?", user.ID).First(&employee).Error
	if empErr != nil && empErr != gorm.ErrRecordNotFound {
		log.Fatal(empErr)
	}

	if empErr == gorm.ErrRecordNotFound {
		employee = models.Employee{
			UserID:    user.ID,
			FirstName: "System",
			LastName:  "Admin",
			Status:    "active",
			HireDate:  time.Now().UTC(),
		}
		if createEmpErr := db.Create(&employee).Error; createEmpErr != nil {
			log.Fatal(createEmpErr)
		}
	}

	fmt.Println("Admin user is ready:")
	fmt.Printf("Email: %s\n", email)
	fmt.Printf("Password: %s\n", password)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

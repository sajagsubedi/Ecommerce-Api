package database

import (
	"fmt"
	"log"
	"os"

	"github.com/sajagsubedi/Ecommerce-Api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() error {
	// Get environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASS")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")
	if sslmode == "" {
		sslmode = "disable" // Default to disable if not set
	}

	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASS", "DB_NAME"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			return fmt.Errorf("missing required environment variable: %s", env)
		}
	}

	// Create the DSN (Data Source Name)
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate all models
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		return fmt.Errorf("error migrating models: %w", err)
	}

	log.Println("Database connection established successfully")
	return nil
}
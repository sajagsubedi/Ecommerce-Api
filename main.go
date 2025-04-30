package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/routes"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error on loading .env file")
	}

	// Initialize database connection
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Get port from environment or default to 8000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Initialize Gin router
	router := gin.New()

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Content-Type",
		"Authorization",
	}

	// Apply middlewares
	router.Use(gin.Logger())
	router.Use(cors.New(config))

	// Set up routes
	routes.AuthRoutes(router)
	routes.ProductRoutes(router)
	routes.CartRoutes(router)
	routes.OrderRoutes(router)

	// Start the server
	log.Printf("Server running on port %s", port)
	router.Run(":" + port)
}

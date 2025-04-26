package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sajagsubedi/Ecommerce-Api/routes"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error on loading .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()

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

	//middlewares
	router.Use(gin.Logger())
	router.Use(cors.New(config))

	//routes
	routes.AuthRoutes(router)

	router.Run(":" + port)

}

package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/models"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()

// HashPassword hashes the provided password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Signup handles user registration.
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Check if database is initialized
		if database.DB == nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Database connection not initialized",
			})
			return
		}

		// Bind and validate input
		user := new(models.User)
		if err := c.BindJSON(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid request payload",
			})
			return
		}

		// Validate struct fields
		if err := validate.Struct(user); err != nil {
			validationErrors := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": validationErrors.Error(),
			})
			return
		}

		// Check if user already exists
		var existingUser []models.User

		// Query the database for a user with the given email
		err := database.DB.Where("email = ?", user.Email).Find(&existingUser).Error
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal Server error",
			})
			return
		}

		//if user already exists with given email
		if len(existingUser) > 0 {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "User with given email already exists!",
			})
			return
		}

		// Hash the password
		hashedPassword, err := HashPassword(user.Password)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to process password",
			})
			return
		}
		user.Password = hashedPassword

		// Create the user
		if err := database.DB.Create(user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to register user",
			})
			return
		}

		// Return success response (exclude sensitive fields)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User registered successfully",
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
			},
		})
	}
}

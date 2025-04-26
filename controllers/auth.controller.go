package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/helpers"
	"github.com/sajagsubedi/Ecommerce-Api/models"
)

var validate = validator.New()

// Signup handles user registration.
func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)

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
		err := db.Where("email = ?", user.Email).Find(&existingUser).Error
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
		hashedPassword, err := user.HashPassword()

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
		if err := db.Create(user).Error; err != nil {
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

func Signin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)

		var user models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal Server Error!",
			})
			return
		}

		if user.Email == "" || user.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Please provide all fields!",
			})
			return
		}

		var existingUser models.User

		err := db.Where("email = ?", user.Email).First(&existingUser).Error

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid credendals!",
			})
			return
		}

		isPasswordCorrect := existingUser.ComparePassword(user.Password)

		if !isPasswordCorrect {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid credendals!",
			})
			return
		}

		token, err := helpers.GenerateToken(string(existingUser.ID))

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Something went wrong!",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logged in successfully!",
			"token":   token,
		})
		return

	}
}

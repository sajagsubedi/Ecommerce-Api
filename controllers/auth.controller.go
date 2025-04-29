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

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)

		user := new(models.User)
		if err := c.BindJSON(user); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := validate.Struct(user); err != nil {
			handleError(c, http.StatusBadRequest, err.(validator.ValidationErrors).Error())
			return
		}

		var existingUser []models.User
		if err := db.Where("email = ?", user.Email).Find(&existingUser).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if len(existingUser) > 0 {
			handleError(c, http.StatusConflict, "User with given email already exists!")
			return
		}

		hashedPassword, err := user.HashPassword()
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to process password")
			return
		}
		user.Password = hashedPassword
		user.Role = "user"

		if err := db.Create(user).Error; err != nil {
			log.Printf("Failed to create user: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to register user")
			return
		}

		cart := models.Cart{
			UserID: user.ID,
			Items:  []models.CartItem{},
		}
		if err := db.Create(&cart).Error; err != nil {
			log.Printf("Failed to create cart: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to create cart")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User registered successfully",
			"user": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
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
			handleError(c, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if user.Email == "" || user.Password == "" {
			handleError(c, http.StatusBadRequest, "Please provide all fields!")
			return
		}

		var existingUser models.User
		if err := db.Where("email = ?", user.Email).First(&existingUser).Error; err != nil {
			handleError(c, http.StatusUnauthorized, "Invalid credentials!")
			return
		}

		if !existingUser.ComparePassword(user.Password) {
			handleError(c, http.StatusUnauthorized, "Invalid credentials!")
			return
		}

		token, err := helpers.GenerateToken(existingUser.ID, existingUser.Role)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Something went wrong!")
			return
		}

		c.SetCookie("Authorization", token, 60*60*24*7, "/", "", false, true)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Logged in successfully!",
			"token":   token,
		})
	}
}

package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/models"
)

func GetProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		fmt.Print(userID)
		var user models.User
		if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Your profile fetched successfully",
			"data": gin.H{
				"id":    user.ID,
				"name":  user.Name,
				"email": user.Email,
				"role":  user.Role,
			},
		})
	}
}

func UpdateProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		type UpdateProfile struct {
			Name string `gorm:"not null" json:"name" validate:"required,min=2,max=100"`
		}

		var input UpdateProfile

		if err := c.ShouldBindJSON(&input); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var existingUser models.User
		if err := db.Where("id = ?", userID).First(&existingUser).Error; err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		if input.Name != "" {
			existingUser.Name = input.Name
		}

		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update user: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to update profile")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Profile updated successfully",
			"user": gin.H{
				"id":    existingUser.ID,
				"name":  existingUser.Name,
				"email": existingUser.Email,
			}})
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		type ChangePasswordRequest struct {
			OldPassword string `json:"old_password" validate:"required"`
			NewPassword string `json:"new_password" validate:"required,min=6"`
		}
		var request ChangePasswordRequest
		if err := c.BindJSON(&request); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var existingUser models.User
		if err := db.Where("id = ?", userID).First(&existingUser).Error; err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		if !existingUser.ComparePassword(request.OldPassword) {
			handleError(c, http.StatusUnauthorized, "Invalid credentials!")
			return
		}

		existingUser.Password = request.NewPassword
		hashedPassword, err := existingUser.HashPassword()
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to process password")
			return
		}
		existingUser.Password = hashedPassword
		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update user: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to update password")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Password changed successfully!",
		})
	}
}

func GetUsersByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)

		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch users")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Users fetched successfully!",
			"data":    users,
		})
	}
}
func GetUserById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		userId := c.Param("userId")

		var user models.User
		err := db.Where("id = ?", userId).First(&user).Error
		if err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User fetched successfully!",
			"data":    user,
		})
	}
}

func UpdateUserByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		userId := c.Param("userId")

		var user models.User
		if err := c.BindJSON(&user); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := validate.Struct(user); err != nil {
			handleError(c, http.StatusBadRequest, err.(validator.ValidationErrors).Error())
			return
		}

		var existingUser models.User
		if err := db.Where("id = ?", userId).First(&existingUser).Error; err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		if user.Email != "" && user.Email != existingUser.Email {
			var emailExists []models.User
			if err := db.Where("email = ?", user.Email).Find(&emailExists).Error; err != nil {
				handleError(c, http.StatusInternalServerError, "Internal Server Error")
				return
			}
			if len(emailExists) > 0 {
				handleError(c, http.StatusConflict, "Email already exists!")
				return
			}
			existingUser.Email = user.Email
		}

		if user.Name != "" {
			existingUser.Name = user.Name
		}

		if user.Password != "" {
			hashedPassword, err := existingUser.HashPassword()
			if err != nil {
				log.Printf("Failed to hash password: %v", err)
				handleError(c, http.StatusInternalServerError, "Failed to process password")
				return
			}
			existingUser.Password = hashedPassword
		}

		if user.Role != "" {
			existingUser.Role = user.Role
		}

		if err := db.Save(&existingUser).Error; err != nil {
			log.Printf("Failed to update user: %v", err)
			handleError(c, http.StatusInternalServerError, "Failed to update user")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User updated successfully",
			"user": gin.H{
				"id":    existingUser.ID,
				"name":  existingUser.Name,
				"email": existingUser.Email,
				"role":  existingUser.Role,
			},
		})
	}
}

func DeleteUserByAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		userId := c.Param("userId")

		var user models.User
		if err := db.Where("id = ?", userId).First(&user).Error; err != nil {
			handleError(c, http.StatusNotFound, "User not found")
			return
		}

		if err := db.Delete(&user).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete user")
			return
		}

		//delete the cart item with the product
		var cart models.Cart
		if err := db.Where("user_id = ?", userId).Delete(&cart).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart not found")
			return
		}

		//delete the cart items
		var cartItems []models.CartItem
		if err := db.Where("cart_id = ?", cart.ID).Delete(&cartItems).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart items not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User deleted successfully",
		})
	}
}

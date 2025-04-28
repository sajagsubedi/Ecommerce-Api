package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/models"
	"gorm.io/gorm"
)

// Reusable error response
func handleError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{"success": false, "message": message})
}

func GetCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		var cart models.Cart
		if err := db.Where("user_id = ?", userID).Preload("Items.Product").First(&cart).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Cart fetched successfully",
			"cart":    cart,
		})
	}
}

func AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		var cartItem models.CartItem
		if err := c.ShouldBindJSON(&cartItem); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := validate.Struct(cartItem); err != nil {
			handleError(c, http.StatusBadRequest, err.Error())
			return
		}

		var cart models.Cart
		if err := db.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart not found")
			return
		}

		var existingItem models.CartItem
		err := db.Where("cart_id = ? AND product_id = ?", cart.ID, cartItem.ProductID).First(&existingItem).Error

		if err == nil {
			existingItem.Quantity += cartItem.Quantity
			if err := db.Save(&existingItem).Error; err != nil {
				handleError(c, http.StatusInternalServerError, "Failed to update cart item")
				return
			}
		} else if err == gorm.ErrRecordNotFound {
			cartItem.CartID = cart.ID
			if err := db.Create(&cartItem).Error; err != nil {
				handleError(c, http.StatusInternalServerError, "Failed to add cart item")
				return
			}
		} else {
			handleError(c, http.StatusInternalServerError, "Database error")
			return
		}

		// Return updated cart
		var updatedCart models.Cart
		if err := db.Where("user_id = ?", userID).Preload("Items.Product").First(&updatedCart).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch updated cart")
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Product added to cart",
			"cart":    updatedCart,
		})
	}
}

func UpdateCartItemQuantity() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cartItemID := c.Param("cartItemId")
		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		var payload struct {
			Quantity int `json:"quantity" binding:"required,min=1"`
		}
		if err := c.ShouldBindJSON(&payload); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		var cartItem models.CartItem
		if err := db.Preload("Product").First(&cartItem, "id = ?", cartItemID).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart item not found")
			return
		}

		var cart models.Cart
		if err := db.First(&cart, "user_id = ?", userID).Error; err != nil || cart.ID != cartItem.CartID {
			handleError(c, http.StatusForbidden, "Unauthorized access to cart item")
			return
		}

		cartItem.Quantity = payload.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update cart item")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Cart item updated",
			"cartItem": cartItem,
		})
	}
}

func DeleteCartItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		cartItemID := c.Param("cartItemId")
		userID, _ := c.Get("userid")
		db := database.DB.WithContext(ctx)

		var cartItem models.CartItem
		if err := db.First(&cartItem, "id = ?", cartItemID).Error; err != nil {
			handleError(c, http.StatusNotFound, "Cart item not found")
			return
		}

		var cart models.Cart
		if err := db.First(&cart, "user_id = ?", userID).Error; err != nil || cart.ID != cartItem.CartID {
			handleError(c, http.StatusForbidden, "Unauthorized access to cart item")
			return
		}

		if err := db.Delete(&cartItem).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete cart item")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Cart item deleted",
		})
	}
}

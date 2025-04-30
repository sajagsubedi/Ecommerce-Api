package controllers

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/models"
	"gorm.io/gorm"
)

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		var products []models.Product

		if err := db.Find(&products).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch products")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Products fetched successfully!",
			"data":    products,
		})
	}
}

func GetProductById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		productId := c.Param("productId")

		var product models.Product
		err := db.Where("id = ?", productId).First(&product).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handleError(c, http.StatusNotFound, "Product not found")
				return
			}
			handleError(c, http.StatusInternalServerError, "Failed to fetch product")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Product fetched successfully!",
			"data":    product,
		})
	}
}

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		var product models.Product

		if err := c.BindJSON(&product); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		if err := db.Create(&product).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to create product")
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Product created successfully!",
			"data":    product,
		})
	}
}

func UpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		productId := c.Param("productId")

		var updatedData models.Product
		if err := c.BindJSON(&updatedData); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// First check if the product exists
		var existingProduct models.Product
		if err := db.Where("id = ?", productId).First(&existingProduct).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handleError(c, http.StatusNotFound, "Product not found")
				return
			}
			handleError(c, http.StatusInternalServerError, "Failed to fetch product")
			return
		}

		// Now update the product
		if err := db.Model(&existingProduct).Updates(updatedData).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update product")
			return
		}

		// Fetch updated product
		if err := db.Where("id = ?", productId).First(&existingProduct).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch updated product")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Product updated successfully!",
			"data":    existingProduct,
		})
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		productId := c.Param("productId")

		var product models.Product
		if err := db.Where("id = ?", productId).First(&product).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				handleError(c, http.StatusNotFound, "Product not found")
				return
			}
			handleError(c, http.StatusInternalServerError, "Failed to find product")
			return
		}

		if err := db.Delete(&product).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete product")
			return
		}

		//delete the cart item with the product
		var cartItem models.CartItem

		if err := db.Where("product_id = ?", productId).Delete(&cartItem).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete cart item using the deleted product")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Product deleted successfully!",
		})
	}
}

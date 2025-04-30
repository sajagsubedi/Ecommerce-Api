package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/models"
)

// OrderItemInput represents the input for an order item
type OrderItemInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

// ShippingAddressInput represents the input for a shipping address
type ShippingAddressInput struct {
	Street  string `json:"street" binding:"required"`
	City    string `json:"city" binding:"required"`
	State   string `json:"state" binding:"required"`
	Country string `json:"country" binding:"required"`
	ZipCode string `json:"zip_code" binding:"required"`
	Notes   string `json:"notes"`
}

// CreateOrderInput represents the input for creating an order
type CreateOrderInput struct {
	ShippingAddress ShippingAddressInput `json:"shipping_address" binding:"required"`
	ContactNumber   string               `json:"contact_number" binding:"required,len=10"`
	Items           []OrderItemInput     `json:"items" binding:"required,dive"`
}

// CreateOrder handles the creation of a new order
func CreateOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, exists := c.Get("userid")
		if !exists {
			handleError(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		var input CreateOrderInput
		if err := c.ShouldBindJSON(&input); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid input data")
			return
		}

		db := database.DB.WithContext(ctx)
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				handleError(c, http.StatusInternalServerError, "Internal server error")
			}
		}()

		// Create order
		order := models.Order{
			ContactNumber: input.ContactNumber,
			Status:        models.OrderStatusPending,
			UserID:        userID.(uint),
			TotalAmount:   0,
		}
		if err := tx.Create(&order).Error; err != nil {
			tx.Rollback()
			handleError(c, http.StatusInternalServerError, "Failed to create order")
			return
		}

		// Create shipping address
		shippingAddress := models.ShippingAddress{
			Street:  input.ShippingAddress.Street,
			City:    input.ShippingAddress.City,
			State:   input.ShippingAddress.State,
			Country: input.ShippingAddress.Country,
			ZipCode: input.ShippingAddress.ZipCode,
			Notes:   input.ShippingAddress.Notes,
			OrderID: order.ID,
		}
		if err := tx.Create(&shippingAddress).Error; err != nil {
			tx.Rollback()
			handleError(c, http.StatusInternalServerError, "Failed to create shipping address")
			return
		}

		// Create order items and update product stock
		var totalAmount float64
		for _, item := range input.Items {
			var product models.Product
			if err := tx.First(&product, item.ProductID).Error; err != nil {
				tx.Rollback()
				handleError(c, http.StatusNotFound, "Product not found")
				return
			}

			if !product.IsAvailable {
				tx.Rollback()
				handleError(c, http.StatusBadRequest, "Product is not available")
				return
			}

			if product.Stock < item.Quantity {
				tx.Rollback()
				handleError(c, http.StatusBadRequest, "Insufficient stock for product")
				return
			}

			product.Stock -= item.Quantity
			if err := tx.Save(&product).Error; err != nil {
				tx.Rollback()
				handleError(c, http.StatusInternalServerError, "Failed to update product stock")
				return
			}

			itemPrice := product.Price * float64(item.Quantity)
			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				Price:     itemPrice,
			}
			if err := tx.Create(&orderItem).Error; err != nil {
				tx.Rollback()
				handleError(c, http.StatusInternalServerError, "Failed to create order item")
				return
			}
			totalAmount += itemPrice
		}

		// Update order with total amount and shipping address
		order.TotalAmount = totalAmount
		if err := tx.Save(&order).Error; err != nil {
			tx.Rollback()
			handleError(c, http.StatusInternalServerError, "Failed to update order")
			return
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			handleError(c, http.StatusInternalServerError, "Failed to commit transaction")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":  true,
			"message":  "Order created successfully",
			"order_id": order.ID,
		})
	}
}

// GetUserOrders retrieves all orders for the authenticated user
func GetUserOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, exists := c.Get("userid")
		if !exists {
			handleError(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		db := database.DB.WithContext(ctx)
		var orders []models.Order
		if err := db.Where("user_id = ?", userID).Preload("ShippingAddress").Preload("Items.Product").Find(&orders).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch orders")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    orders,
			"message": "Orders retrieved successfully",
		})
	}
}

// GetUserOrderByID retrieves a specific order by ID for the authenticated user
func GetUserOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, exists := c.Get("userid")
		if !exists {
			handleError(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		orderID := c.Param("id")
		db := database.DB.WithContext(ctx)
		var order models.Order
		if err := db.Where("user_id = ? AND id = ?", userID, orderID).Preload("ShippingAddress").Preload("Items.Product").First(&order).Error; err != nil {
			handleError(c, http.StatusNotFound, "Order not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    order,
			"message": "Order retrieved successfully",
		})
	}
}

// CancelUserOrder cancels a pending order for the authenticated user
func CancelUserOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID, exists := c.Get("userid")
		if !exists {
			handleError(c, http.StatusUnauthorized, "User not authenticated")
			return
		}

		orderID := c.Param("id")
		db := database.DB.WithContext(ctx)
		var order models.Order
		if err := db.Where("user_id = ? AND id = ?", userID, orderID).First(&order).Error; err != nil {
			handleError(c, http.StatusNotFound, "Order not found")
			return
		}

		if order.Status != models.OrderStatusPending {
			handleError(c, http.StatusBadRequest, "Only pending orders can be cancelled")
			return
		}

		order.Status = models.OrderStatusCancelled
		if err := db.Save(&order).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to cancel order")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Order cancelled successfully",
		})
	}
}

// AdminGetAllOrders retrieves all orders for admin users
func AdminGetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		var orders []models.Order
		if err := db.Preload("ShippingAddress").Preload("Items.Product").Find(&orders).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch orders")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    orders,
			"message": "Orders retrieved successfully",
		})
	}
}

// AdminGetOrderByID retrieves a specific order by ID for admin users
func AdminGetOrderByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderID := c.Param("id")
		db := database.DB.WithContext(ctx)
		var order models.Order
		if err := db.Where("id = ?", orderID).Preload("ShippingAddress").Preload("Items.Product").First(&order).Error; err != nil {
			handleError(c, http.StatusNotFound, "Order not found")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    order,
			"message": "Order retrieved successfully",
		})
	}
}

// AdminGetOrdersByUserID retrieves all orders for a specific user for admin users
func AdminGetOrdersByUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		userID := c.Param("user_id")
		db := database.DB.WithContext(ctx)
		var orders []models.Order
		if err := db.Where("user_id = ?", userID).Preload("ShippingAddress").Preload("Items.Product").Find(&orders).Error; err != nil {
			handleError(c, http.StatusNotFound, "Orders not found for this user")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    orders,
			"message": "Orders retrieved successfully",
		})
	}
}

// AdminGetAllOrderItems retrieves all order items for admin users
func AdminGetAllOrderItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		db := database.DB.WithContext(ctx)
		var orderItems []models.OrderItem
		if err := db.Preload("Product").Find(&orderItems).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to fetch order items")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    orderItems,
			"message": "Order items retrieved successfully",
		})
	}
}

// AdminGetOrderItemsByProductID retrieves order items by product ID for admin users
func AdminGetOrderItemsByProductID() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		productID := c.Param("product_id")
		db := database.DB.WithContext(ctx)
		var orderItems []models.OrderItem
		if err := db.Where("product_id = ?", productID).Preload("Product").Find(&orderItems).Error; err != nil {
			handleError(c, http.StatusNotFound, "Order items not found for this product")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    orderItems,
			"message": "Order items retrieved successfully",
		})
	}
}

// AdminUpdateOrderStatus updates the status of an order for admin users
func AdminUpdateOrderStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		orderID := c.Param("id")
		db := database.DB.WithContext(ctx)
		var order models.Order
		if err := db.Where("id = ?", orderID).First(&order).Error; err != nil {
			handleError(c, http.StatusNotFound, "Order not found")
			return
		}

		var input struct {
			Status string `json:"status" binding:"required,oneof=pending processing shipped delivered cancelled"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			handleError(c, http.StatusBadRequest, "Invalid input data")
			return
		}

		order.Status = input.Status
		if err := db.Save(&order).Error; err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update order status")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Order status updated successfully",
		})
	}
}

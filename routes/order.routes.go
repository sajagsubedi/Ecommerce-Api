package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sajagsubedi/Ecommerce-Api/controllers"
	"github.com/sajagsubedi/Ecommerce-Api/middlewares"
)

// OrderRoutes sets up the API routes for order management
func OrderRoutes(incomingRoutes *gin.Engine) {
	// User order routes
	userOrderRoutes := incomingRoutes.Group("/api/v1/orders")
	userOrderRoutes.Use(middlewares.CheckUser())
	userOrderRoutes.POST("/checkout", controllers.CreateOrder())
	userOrderRoutes.GET("/", controllers.GetUserOrders())
	userOrderRoutes.GET("/:id", controllers.GetUserOrderByID())
	userOrderRoutes.DELETE("/:id/cancel", controllers.CancelUserOrder())

	// Admin order routes
	adminOrderRoutes := incomingRoutes.Group("/api/v1/admin/orders")
	adminOrderRoutes.Use(middlewares.CheckAdmin())
	adminOrderRoutes.GET("/", controllers.AdminGetAllOrders())
	adminOrderRoutes.GET("/:id", controllers.AdminGetOrderByID())
	adminOrderRoutes.GET("/user/:user_id", controllers.AdminGetOrdersByUserID())
	adminOrderRoutes.PUT("/:id/status", controllers.AdminUpdateOrderStatus())

	// Admin order item routes
	adminOrderItemRoutes := incomingRoutes.Group("/api/v1/admin/order-items")
	adminOrderItemRoutes.Use(middlewares.CheckAdmin())
	adminOrderItemRoutes.GET("/", controllers.AdminGetAllOrderItems())
	adminOrderItemRoutes.GET("/product/:product_id", controllers.AdminGetOrderItemsByProductID())
}

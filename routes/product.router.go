package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
	"github.com/sajagsubedi/Ecommerce-Api/middlewares"
)

func ProductRoutes(incomingRoutes *gin.Engine) {
	productRoutes := incomingRoutes.Group("/api/v1/products")
	productRoutes.GET("/", controller.GetAllProducts())
	productRoutes.GET("/:productId", controller.GetProductById())

	adminRoutes := productRoutes.Group("")
	adminRoutes.Use(middlewares.CheckAdmin())
	adminRoutes.POST("/", controller.CreateProduct())
	productRoutes.PUT("/:productId", controller.UpdateProduct())
	productRoutes.DELETE("/:productId", controller.DeleteProduct())
}

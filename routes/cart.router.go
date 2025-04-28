package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
	"github.com/sajagsubedi/Ecommerce-Api/middlewares"
)

func CartRoutes(incomingRoutes *gin.Engine) {
	CartRoutes := incomingRoutes.Group("/api/v1/cart")
	CartRoutes.Use(middlewares.CheckUser())
	CartRoutes.GET("/", controller.GetCart())
	CartRoutes.POST("/", controller.AddToCart())
	CartRoutes.PUT("/update-quantity/:cartItemId", controller.UpdateCartItemQuantity())
	CartRoutes.DELETE("/:id", controller.DeleteCartItem())
}

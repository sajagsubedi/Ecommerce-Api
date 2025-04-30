package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
	"github.com/sajagsubedi/Ecommerce-Api/middlewares"
)

func CartRoutes(incomingRoutes *gin.Engine) {
	incomingcartRoutes := incomingRoutes.Group("/api/v1/cart")
	incomingcartRoutes.Use(middlewares.CheckUser())
	incomingcartRoutes.GET("/", controller.GetCart())
	incomingcartRoutes.POST("/", controller.AddToCart())
	incomingcartRoutes.PUT("/update-quantity/:cartItemId", controller.UpdateCartItemQuantity())
	incomingcartRoutes.DELETE("/:id", controller.DeleteCartItem())
}

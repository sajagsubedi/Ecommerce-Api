package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	authRoutes := incomingRoutes.Group("/api/v1/auth")
	authRoutes.POST("/signup", controller.Signup())
	authRoutes.POST("/signin", controller.Signin())
}

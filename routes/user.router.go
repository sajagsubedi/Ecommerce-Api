package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
	"github.com/sajagsubedi/Ecommerce-Api/middlewares"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	authRoutes := incomingRoutes.Group("/api/v1/user")
	authRoutes.Use(middlewares.CheckUser())
	authRoutes.GET("/profile", controller.GetProfile())
	authRoutes.PUT("/profile", controller.UpdateProfile())
	authRoutes.POST("/change-password", controller.ChangePassword())

	adminRoutes := incomingRoutes.Group("/api/v1/admin/users")
	adminRoutes.Use(middlewares.CheckAdmin())

	adminRoutes.GET("/", controller.GetUsersByAdmin())
	adminRoutes.GET("/:userId", controller.GetUserById())
	adminRoutes.PUT("/:userId", controller.UpdateUserByAdmin())
	adminRoutes.DELETE("/:userId", controller.DeleteUserByAdmin())
}

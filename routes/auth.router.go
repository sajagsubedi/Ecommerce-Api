package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/sajagsubedi/Ecommerce-Api/controllers"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/signup", controller.Signup())
}

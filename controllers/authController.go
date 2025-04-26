package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Signup Successfull!",
		})
	}
}

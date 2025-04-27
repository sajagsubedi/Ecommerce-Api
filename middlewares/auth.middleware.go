package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sajagsubedi/Ecommerce-Api/database"
	"github.com/sajagsubedi/Ecommerce-Api/helpers"
	"github.com/sajagsubedi/Ecommerce-Api/models"
)

func DecodeJwt(c *gin.Context) (*models.User, string) {
	authToken := c.Request.Header.Get("Authorization")

	if authToken == "" {
		cookieToken, err := c.Cookie("Authorization")
		if err != nil {
			msg := "No Authorization header or cookie provided"
			return nil, msg
		}
		authToken = cookieToken
	} else {
		if strings.HasPrefix(authToken, "Bearer ") {
			authToken = strings.TrimPrefix(authToken, "Bearer ")
		} else {
			return nil, "Invalid token format. Include \"Brearer \" prefix "
		}
	}

	//get the claims from the token
	claims, err := helpers.ValidateToken(authToken)

	if err != nil {
		return nil, err.Error()
	}

	existingUser := new(models.User)

	err = database.DB.Where("id = ?", claims.UserID).First(existingUser).Error

	if err != nil {
		return nil, err.Error()
	}

	fmt.Print(existingUser)
	return existingUser, ""
}

func CheckUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, msg := DecodeJwt(c)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": msg,
			})
			c.Abort()
			return
		}

		c.Set("userid", claims.ID)
		c.Set("usertype", claims.Role)
		c.Next()
	}
}

func CheckAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, msg := DecodeJwt(c)
		if msg != "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": msg,
			})
			c.Abort()
			return
		}

		if claims.Role != "admin" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false, "message": "Unauthorized access",
			})
			c.Abort()
			return
		}
		c.Set("userid", claims.ID)
		c.Set("usertype", claims.Role)
		c.Next()
	}
}

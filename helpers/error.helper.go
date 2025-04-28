package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HandleInternalServerError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Internal Server Error",
		})
		return true
	}
	return false
}

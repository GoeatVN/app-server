package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware to handle errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Error: %s\n", e.Error())
				if e.Type == gin.ErrorTypePublic {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": e.Error(),
					})
					return
				}
				if e.Type == gin.ErrorTypeBind {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": e.Error(),
					})
					return
				}

			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Server Error",
			})
		}

	}
}

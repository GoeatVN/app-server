package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleError middleware to handle errors
func ErrorHandler(c *gin.Context) {
	// Process request
	// Check if there are any errors
	if len(c.Errors) > 0 {
		// Log the error
		for _, e := range c.Errors {
			// You can log the error here
			// log.Error(e.Err)
			// For simplicity, we just print it
			println(e.Err.Error())
		}
		// Return a generic error message to the client
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal Server Error",
		})
		return
	}
	c.Next()
}

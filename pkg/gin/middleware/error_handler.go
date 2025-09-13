package middleware

import (
	"net/http"

	apperr "github.com/busnosh/go-utils/pkg/error"
	"github.com/busnosh/go-utils/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ErrorHandler logs errors and returns a structured JSON response
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process request

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err

			// Default response
			status := http.StatusInternalServerError
			message := "Internal Server Error"

			// If it's an AppError, use its code & message
			if appErr, ok := lastErr.(*apperr.Error); ok {
				status = appErr.Code
				message = appErr.Message
			}

			// Log the error using your custom logger
			logger.Error("Request %s %s -> %d | Error: %v", c.Request.Method, c.Request.URL.Path, status, lastErr)

			// Respond with structured JSON
			c.JSON(status, gin.H{
				"success": false,
				"message": message,
			})
		}
	}
}

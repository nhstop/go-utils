package middleware

import (
	"net/http"

	"github.com/busnosh/go-utils/pkg/constants"
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

			statusColor := constants.ColorGreen
			if status >= 400 && status < 500 {
				statusColor = constants.ColorYellow
			} else if status >= 500 {
				statusColor = constants.ColorRed
			}

			logger.Error("%sRequest %s %s -> %s%d%s | Error: %s%v%s",
				constants.ColorBlue, // HTTP method + path in blue
				c.Request.Method,
				c.Request.URL.Path,
				statusColor, status, constants.ColorReset, // status in green/yellow/red
				constants.ColorRed, lastErr, constants.ColorReset, // error message in red
			)

			// Respond with structured JSON
			c.JSON(status, gin.H{
				"success": false,
				"message": message,
			})
		}
	}
}

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
		c.Next()

		if len(c.Errors) > 0 {
			lastErr := c.Errors.Last().Err

			// Default response
			status := http.StatusInternalServerError
			message := "Internal Server Error"
			code := 0

			// If it's an AppError, use its HTTPCode and message
			if appErr, ok := lastErr.(*apperr.Error); ok {
				status = appErr.HTTPCode // ✅ use HTTPCode
				message = appErr.Message
				code = appErr.Code
			}

			statusColor := constants.ColorGreen
			if status >= 400 && status < 500 {
				statusColor = constants.ColorYellow
			} else if status >= 500 {
				statusColor = constants.ColorRed
			}
			logger.Error("%sRequest %s %s -> %s%d%s | Code: %s%d%s | Error: %s%v%s",
				constants.ColorBlue,
				c.Request.Method,
				c.Request.URL.Path,
				statusColor, status, constants.ColorReset,
				constants.ColorYellow, code, constants.ColorReset,
				constants.ColorRed, lastErr, constants.ColorReset,
			)

			// Respond with structured JSON
			c.JSON(status, gin.H{
				"success": false,
				"message": message,
				"code":    code,
			})

		}
	}
}

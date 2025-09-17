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

		if len(c.Errors) == 0 {
			return
		}

		lastErr := c.Errors.Last().Err

		// Default response
		status := http.StatusInternalServerError
		message := "Internal Server Error"
		code := constants.Empty

		// If it's a CodedError, use its HTTPCode, Message, and Code
		if codedErr, ok := lastErr.(*apperr.CodedError); ok {
			status = codedErr.HTTPCode
			message = codedErr.Message
			code = codedErr.Code
		}

		// Color coding for logs
		statusColor := constants.ColorGreen
		if status >= 400 && status < 500 {
			statusColor = constants.ColorYellow
		} else if status >= 500 {
			statusColor = constants.ColorRed
		}

		// Log the error
		logger.Error("%sRequest %s %s -> %s%d%s | Code: %s%d%s | Error: %s%v%s",
			constants.ColorBlue,
			c.Request.Method,
			c.Request.URL.Path,
			statusColor, status, constants.ColorReset,
			constants.ColorYellow, code, constants.ColorReset,
			constants.ColorRed, lastErr, constants.ColorReset,
		)

		// Respond with structured JSON
		resp := gin.H{
			"success": false,
			"message": message,
		}

		if code != 0 {
			resp["code"] = code
		}

		c.JSON(status, resp)

	}
}

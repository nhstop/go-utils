package middlewares

import (
	"time"

	"github.com/busnosh/go-utils/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Choose status color
		statusColor := ColorGreen
		switch {
		case status >= 500:
			statusColor = ColorRed
		case status >= 400:
			statusColor = ColorYellow
		case status >= 300:
			statusColor = ColorBlue
		}

		// Capture errors if any
		errorMessage := ""
		if len(c.Errors) > 0 {
			errorMessage = c.Errors.String()
		}

		// Log the request using your custom logger
		logger.Info("%s %s | %s%d%s | %v", method, path, statusColor, status, ColorReset, latency)
		if errorMessage != "" {
			logger.Error("Error: %s%s%s", statusColor, errorMessage, ColorReset)
		}
	}
}

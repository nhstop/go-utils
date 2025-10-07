package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nhstop/go-utils/pkg/constants"
	"github.com/nhstop/go-utils/pkg/logger"
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
		statusColor := constants.ColorGreen
		switch {
		case status >= 500:
			statusColor = constants.ColorRed
		case status >= 400:
			statusColor = constants.ColorYellow
		case status >= 300:
			statusColor = constants.ColorBlue
		}

		// Log the request using your custom logger
		logger.Info("%s %s | %s%d%s | %v", method, path, statusColor, status, constants.ColorReset, latency)

	}
}

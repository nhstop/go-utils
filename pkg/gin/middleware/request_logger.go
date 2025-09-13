package middleware

import (
	"time"

	"github.com/busnosh/go-utils/pkg/log"
	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// After request
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		// Use your custom logger
		log.Info("%s %s -> %d (%v)", method, path, status, duration)
	}
}

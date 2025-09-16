package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig allows customization of headers and CORS
type SecurityHeadersConfig struct {
	AllowedOrigins []string // e.g. []string{"http://localhost:3000", "https://example.com"}
	AllowedMethods []string // e.g. GET, POST, PUT, DELETE
	AllowedHeaders []string // e.g. Origin, Content-Type, Authorization
}

// DefaultConfig returns a safe default configuration
func DefaultConfig() *SecurityHeadersConfig {
	return &SecurityHeadersConfig{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}
}

// Features included
// Security headers: X-Content-Type-Options, X-Frame-Options, X-XSS-Protection
// Referrer policy: Protect privacy
// CSP: Prevents XSS and injection
// HSTS: Enforces HTTPS
// COOP + COEP: Cross-origin isolation
// Cache-Control: Protect sensitive pages
// Expect-CT: Certificate transparency enforcement
// CORS: Configurable origins, methods, headers, and preflight handling
// Remove Server header: Reduces fingerprinting
// Optional secure cookie example

// SecurityHeaders middleware adds security headers + CORS + removes Server header
func SecurityHeaders(cfg *SecurityHeadersConfig) gin.HandlerFunc {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	return func(c *gin.Context) {
		// ---------------- Security Headers ----------------
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Content-Security-Policy", "default-src 'self'; img-src 'self' data:; script-src 'self'; style-src 'self' 'unsafe-inline'; font-src 'self'; connect-src 'self'")
		c.Header("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		c.Header("Cross-Origin-Opener-Policy", "same-origin")
		c.Header("Cross-Origin-Embedder-Policy", "require-corp")
		c.Header("Cache-Control", "no-store, no-cache, must-revalidate, private")
		c.Header("Expect-CT", "max-age=86400, enforce")

		// ---------------- CORS ----------------
		origin := c.GetHeader("Origin")
		for _, o := range cfg.AllowedOrigins {
			if o == "*" || o == origin {
				c.Header("Access-Control-Allow-Origin", o)
				c.Header("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
				c.Header("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
				c.Header("Access-Control-Allow-Credentials", "true")
				break
			}
		}

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		// ---------------- Remove Server header ----------------
		c.Writer.Header().Del("Server")

		// ---------------- Optional: Secure cookies example ----------------
		// Uncomment and use when setting cookies:
		// http.SetCookie(c.Writer, &http.Cookie{
		//     Name:     "session",
		//     Value:    "xyz",
		//     HttpOnly: true,
		//     Secure:   true,
		//     SameSite: http.SameSiteStrictMode,
		// })

		c.Next()
	}
}

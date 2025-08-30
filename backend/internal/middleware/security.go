package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// SecurityHeaders adds security-related HTTP headers
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		
		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' ws: wss:; frame-ancestors 'none';")
		
		// Strict Transport Security (only for HTTPS)
		if c.Request.TLS != nil {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}
		
		c.Next()
	}
}

// RateLimit implements rate limiting using token bucket algorithm
func RateLimit(limit rate.Limit) gin.HandlerFunc {
	// Create a limiter for each IP address
	limiters := make(map[string]*rate.Limiter)
	
	return func(c *gin.Context) {
		// Get client IP
		clientIP := getClientIP(c)
		
		// Get or create limiter for this IP
		limiter, exists := limiters[clientIP]
		if !exists {
			limiter = rate.NewLimiter(limit, int(limit))
			limiters[clientIP] = limiter
		}
		
		// Check if request is allowed
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":   "Rate limit exceeded",
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests. Please try again later.",
				"retry_after": time.Now().Add(time.Second).Unix(),
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID is already set
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a new request ID
			requestID = generateRequestID()
		}
		
		// Set request ID in context and response headers
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		
		c.Next()
	}
}

// Logging middleware for request/response logging
func Logging() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("[%s] %s %s %d %s %s %s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.ClientIP,
			param.ErrorMessage,
		)
	})
}

// Recovery middleware for panic recovery
func Recovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic
		log.Printf("Panic recovered: %v", recovered)
		
		// Return internal server error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Internal server error",
			"code":    "INTERNAL_ERROR",
			"message": "An unexpected error occurred. Please try again later.",
		})
	})
}

// getClientIP gets the real client IP address
func getClientIP(c *gin.Context) string {
	// Check for forwarded headers
	if forwardedFor := c.GetHeader("X-Forwarded-For"); forwardedFor != "" {
		return forwardedFor
	}
	
	if realIP := c.GetHeader("X-Real-IP"); realIP != "" {
		return realIP
	}
	
	// Fall back to remote address
	return c.ClientIP()
}

// generateRequestID generates a unique request ID
func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
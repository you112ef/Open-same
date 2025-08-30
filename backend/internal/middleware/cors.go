package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS middleware handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	config := cors.DefaultConfig()
	
	// Allow all origins in development, restrict in production
	config.AllowAllOrigins = true
	
	// Allow specific origins in production
	// config.AllowedOrigins = []string{
	//     "https://yourdomain.com",
	//     "https://www.yourdomain.com",
	//     "https://app.yourdomain.com",
	// }
	
	// Allow specific methods
	config.AllowMethods = []string{
		"GET",
		"POST",
		"PUT",
		"PATCH",
		"DELETE",
		"HEAD",
		"OPTIONS",
	}
	
	// Allow specific headers
	config.AllowHeaders = []string{
		"Origin",
		"Content-Length",
		"Content-Type",
		"Authorization",
		"Accept",
		"Accept-Encoding",
		"Accept-Language",
		"Cache-Control",
		"Connection",
		"DNT",
		"Host",
		"Pragma",
		"Referer",
		"User-Agent",
		"X-Requested-With",
		"X-Forwarded-For",
		"X-Forwarded-Proto",
		"X-Real-IP",
	}
	
	// Allow credentials (cookies, authorization headers)
	config.AllowCredentials = true
	
	// Expose headers to the client
	config.ExposeHeaders = []string{
		"Content-Length",
		"Content-Type",
		"Content-Disposition",
		"X-Total-Count",
		"X-Page-Count",
		"X-Current-Page",
		"X-Per-Page",
		"X-Request-ID",
		"X-Response-Time",
	}
	
	// Set max age for preflight requests
	config.MaxAge = 86400 // 24 hours
	
	return cors.New(config)
}
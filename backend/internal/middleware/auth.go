package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/open-same/backend/internal/config"
	"github.com/open-same/backend/internal/database"
	"github.com/open-same/backend/internal/models"
)

// Claims represents JWT claims
type Claims struct {
	UserID   string `json:"user_id"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsAdmin  bool   `json:"is_admin"`
	jwt.RegisteredClaims
}

// Auth middleware validates JWT tokens and sets user context
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Authorization header required",
				"code":    "MISSING_AUTH_HEADER",
				"message": "Please provide a valid authorization token",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header format",
				"code":    "INVALID_AUTH_FORMAT",
				"message": "Authorization header must start with 'Bearer '",
			})
			c.Abort()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			var errorMessage string
			var errorCode string

			if strings.Contains(err.Error(), "token is expired") {
				errorMessage = "Token has expired"
				errorCode = "TOKEN_EXPIRED"
			} else if strings.Contains(err.Error(), "signature is invalid") {
				errorMessage = "Invalid token signature"
				errorCode = "INVALID_SIGNATURE"
			} else {
				errorMessage = "Invalid token"
				errorCode = "INVALID_TOKEN"
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   errorMessage,
				"code":    errorCode,
				"message": "Please provide a valid authorization token",
			})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token claims",
				"code":    "INVALID_CLAIMS",
				"message": "Token contains invalid claims",
			})
			c.Abort()
			return
		}

		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Token has expired",
				"code":    "TOKEN_EXPIRED",
				"message": "Please refresh your token",
			})
			c.Abort()
			return
		}

		// Get user from database
		var user models.User
		userID, err := parseUUID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid user ID in token",
				"code":    "INVALID_USER_ID",
				"message": "Token contains invalid user information",
			})
			c.Abort()
			return
		}

		if err := database.GetDB().First(&user, "id = ?", userID).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "User not found",
				"code":    "USER_NOT_FOUND",
				"message": "User associated with token not found",
			})
			c.Abort()
			return
		}

		// Check if user is active
		if !user.IsActive {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "User account is deactivated",
				"code":    "USER_DEACTIVATED",
				"message": "Your account has been deactivated",
			})
			c.Abort()
			return
		}

		// Set user context
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Set("is_admin", user.IsAdmin)

		c.Next()
	}
}

// AdminOnly middleware ensures only admin users can access
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		isAdmin, exists := c.Get("is_admin")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "User context not found",
				"code":    "MISSING_USER_CONTEXT",
				"message": "Internal server error",
			})
			c.Abort()
			return
		}

		if !isAdmin.(bool) {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Admin access required",
				"code":    "ADMIN_ACCESS_REQUIRED",
				"message": "You don't have permission to access this resource",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth middleware provides optional authentication
func OptionalAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No token provided, continue without authentication
			c.Next()
			return
		}

		// Check if header starts with "Bearer "
		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Claims)
		if !ok || !token.Valid {
			// Invalid claims, continue without authentication
			c.Next()
			return
		}

		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			// Expired token, continue without authentication
			c.Next()
			return
		}

		// Get user from database
		var user models.User
		userID, err := parseUUID(claims.UserID)
		if err != nil {
			// Invalid user ID, continue without authentication
			c.Next()
			return
		}

		if err := database.GetDB().First(&user, "id = ?", userID).Error; err != nil {
			// User not found, continue without authentication
			c.Next()
			return
		}

		// Check if user is active
		if !user.IsActive {
			// User deactivated, continue without authentication
			c.Next()
			return
		}

		// Set user context
		c.Set("user", &user)
		c.Set("user_id", user.ID)
		c.Set("is_admin", user.IsAdmin)

		c.Next()
	}
}

// GetUserFromContext gets the authenticated user from context
func GetUserFromContext(c *gin.Context) (*models.User, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}
	return user.(*models.User), true
}

// GetUserIDFromContext gets the authenticated user ID from context
func GetUserIDFromContext(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// IsAdmin checks if the authenticated user is an admin
func IsAdmin(c *gin.Context) bool {
	isAdmin, exists := c.Get("is_admin")
	if !exists {
		return false
	}
	return isAdmin.(bool)
}

// parseUUID parses a UUID string
func parseUUID(s string) (string, error) {
	// Simple validation - in production you might want to use a proper UUID library
	if len(s) != 36 {
		return "", fmt.Errorf("invalid UUID format")
	}
	return s, nil
}
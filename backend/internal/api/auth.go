package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/open-same/backend/internal/config"
	"github.com/open-same/backend/internal/database"
	"github.com/open-same/backend/internal/middleware"
	"github.com/open-same/backend/internal/models"
)

// AuthRequest represents authentication request
type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Username  string `json:"username" binding:"required,min=3,max=30"`
	Password  string `json:"password" binding:"required,min=6"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	TokenType    string      `json:"token_type"`
	ExpiresIn    int64       `json:"expires_in"`
	User         models.User `json:"user"`
}

// RefreshRequest represents token refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register handles user registration
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := database.GetDB().Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "User already exists",
			"code":    "USER_EXISTS",
			"message": "A user with this email or username already exists",
		})
		return
	}

	// Create new user
	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	// Set password
	if err := user.SetPassword(req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"code":    "PASSWORD_HASH_ERROR",
			"message": "An error occurred while creating your account",
		})
		return
	}

	// Save user to database
	if err := database.GetDB().Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create user",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while creating your account",
		})
		return
	}

	// Generate tokens
	cfg := config.Load()
	accessToken, refreshToken, err := generateTokens(&user, cfg.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate tokens",
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "An error occurred while creating your account",
		})
		return
	}

	// Save refresh token to database
	token := models.Token{
		UserID:    user.ID,
		Token:     refreshToken,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(time.Duration(cfg.JWT.RefreshHours) * time.Hour),
	}

	if err := database.GetDB().Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save token",
			"code":    "TOKEN_SAVE_ERROR",
			"message": "An error occurred while creating your account",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"data": AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(cfg.JWT.ExpirationHours * 3600),
			User:         user,
		},
	})
}

// Login handles user authentication
func Login(c *gin.Context) {
	var req AuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// Find user by email
	var user models.User
	if err := database.GetDB().Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"code":    "INVALID_CREDENTIALS",
			"message": "Email or password is incorrect",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Account deactivated",
			"code":    "ACCOUNT_DEACTIVATED",
			"message": "Your account has been deactivated",
		})
		return
	}

	// Verify password
	if !user.CheckPassword(req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"code":    "INVALID_CREDENTIALS",
			"message": "Email or password is incorrect",
		})
		return
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	database.GetDB().Save(&user)

	// Generate tokens
	cfg := config.Load()
	accessToken, refreshToken, err := generateTokens(&user, cfg.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate tokens",
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "An error occurred while logging in",
		})
		return
	}

	// Save refresh token to database
	token := models.Token{
		UserID:    user.ID,
		Token:     refreshToken,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(time.Duration(cfg.JWT.RefreshHours) * time.Hour),
	}

	if err := database.GetDB().Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save token",
			"code":    "TOKEN_SAVE_ERROR",
			"message": "An error occurred while logging in",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(cfg.JWT.ExpirationHours * 3600),
			User:         user,
		},
	})
}

// RefreshToken handles token refresh
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// Find refresh token in database
	var token models.Token
	if err := database.GetDB().Where("token = ? AND type = ? AND is_revoked = ?", req.RefreshToken, "refresh", false).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid refresh token",
			"code":    "INVALID_REFRESH_TOKEN",
			"message": "Invalid or expired refresh token",
		})
		return
	}

	// Check if token is expired
	if token.IsExpired() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Refresh token expired",
			"code":    "REFRESH_TOKEN_EXPIRED",
			"message": "Refresh token has expired",
		})
		return
	}

	// Get user
	var user models.User
	if err := database.GetDB().First(&user, token.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "User not found",
			"code":    "USER_NOT_FOUND",
			"message": "User associated with token not found",
		})
		return
	}

	// Check if user is active
	if !user.IsActive {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Account deactivated",
			"code":    "ACCOUNT_DEACTIVATED",
			"message": "Your account has been deactivated",
		})
		return
	}

	// Revoke old refresh token
	token.Revoke()
	database.GetDB().Save(&token)

	// Generate new tokens
	cfg := config.Load()
	accessToken, refreshToken, err := generateTokens(&user, cfg.JWT)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate tokens",
			"code":    "TOKEN_GENERATION_ERROR",
			"message": "An error occurred while refreshing tokens",
		})
		return
	}

	// Save new refresh token
	newToken := models.Token{
		UserID:    user.ID,
		Token:     refreshToken,
		Type:      "refresh",
		ExpiresAt: time.Now().Add(time.Duration(cfg.JWT.RefreshHours) * time.Hour),
	}

	if err := database.GetDB().Create(&newToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save token",
			"code":    "TOKEN_SAVE_ERROR",
			"message": "An error occurred while refreshing tokens",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Tokens refreshed successfully",
		"data": AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
			ExpiresIn:    int64(cfg.JWT.ExpirationHours * 3600),
			User:         user,
		},
	})
}

// generateTokens generates access and refresh tokens
func generateTokens(user *models.User, jwtConfig config.JWTConfig) (string, string, error) {
	// Generate access token
	accessClaims := middleware.Claims{
		UserID:   user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.ExpirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "open-same",
			Subject:   user.ID.String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshClaims := middleware.Claims{
		UserID:   user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		IsAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(jwtConfig.RefreshHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "open-same",
			Subject:   user.ID.String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(jwtConfig.Secret))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
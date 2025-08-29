package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/open-same/backend/internal/api"
	"github.com/open-same/backend/internal/config"
	"github.com/open-same/backend/internal/database"
	"github.com/open-same/backend/internal/middleware"
	"github.com/open-same/backend/internal/redis"
	"github.com/open-same/backend/internal/websocket"
	"golang.org/x/time/rate"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Init(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Redis
	redisClient, err := redis.Init(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Initialize WebSocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create router
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RateLimit(rate.Limit(cfg.RateLimit)))
	router.Use(middleware.RequestID())
	router.Use(middleware.SecurityHeaders())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "healthy",
			"timestamp": time.Now().UTC(),
			"version":   cfg.Version,
		})
	})

	// API routes
	apiGroup := router.Group("/api/v1")
	{
		// Public routes
		apiGroup.GET("/docs", api.ServeDocs)
		apiGroup.POST("/auth/register", api.Register)
		apiGroup.POST("/auth/login", api.Login)
		apiGroup.POST("/auth/refresh", api.RefreshToken)
		apiGroup.GET("/content/public", api.GetPublicContent)

		// Protected routes
		protected := apiGroup.Group("/")
		protected.Use(middleware.Auth(cfg.JWT.Secret))
		{
			// User management
			protected.GET("/user/profile", api.GetUserProfile)
			protected.PUT("/user/profile", api.UpdateUserProfile)
			protected.DELETE("/user/account", api.DeleteUserAccount)

			// Content management
			protected.POST("/content", api.CreateContent)
			protected.GET("/content", api.GetUserContent)
			protected.GET("/content/:id", api.GetContent)
			protected.PUT("/content/:id", api.UpdateContent)
			protected.DELETE("/content/:id", api.DeleteContent)
			protected.POST("/content/:id/share", api.ShareContent)
			protected.POST("/content/:id/collaborate", api.AddCollaborator)

			// Collaboration
			protected.GET("/collaborations", api.GetCollaborations)
			protected.PUT("/collaborations/:id", api.UpdateCollaboration)
			protected.DELETE("/collaborations/:id", api.RemoveCollaborator)

			// Real-time collaboration
			protected.GET("/ws", func(c *gin.Context) {
				websocket.HandleWebSocket(wsHub, c.Writer, c.Request)
			})
		}

		// Admin routes
		admin := apiGroup.Group("/admin")
		admin.Use(middleware.AdminOnly())
		{
			admin.GET("/users", api.AdminGetUsers)
			admin.GET("/content", api.AdminGetAllContent)
			admin.GET("/stats", api.AdminGetStats)
			admin.POST("/users/:id/ban", api.AdminBanUser)
		}
	}

	// GraphQL endpoint
	router.POST("/graphql", api.GraphQLHandler)

	// WebSocket endpoint for real-time collaboration
	router.GET("/ws", func(c *gin.Context) {
		websocket.HandleWebSocket(wsHub, c.Writer, c.Request)
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}
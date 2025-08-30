package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Environment string
	Version     string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	RabbitMQ    RabbitMQConfig
	JWT         JWTConfig
	AI          AIConfig
	RateLimit   float64
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port         int
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database connection configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
	SSLMode  string
}

// RedisConfig holds Redis connection configuration
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// RabbitMQConfig holds RabbitMQ connection configuration
type RabbitMQConfig struct {
	Host     string
	Port     int
	User     string
	Password string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret           string
	ExpirationHours int
	RefreshHours     int
}

// AIConfig holds AI service configuration
type AIConfig struct {
	OpenAIKey      string
	OpenAIModel    string
	AnthropicKey   string
	AnthropicModel string
	MaxTokens      int
	Temperature    float64
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Environment: getEnv("ENVIRONMENT", "development"),
		Version:     getEnv("VERSION", "1.0.0"),
		Server: ServerConfig{
			Port:         getEnvAsInt("API_PORT", 8080),
			Host:         getEnv("API_HOST", "0.0.0.0"),
			ReadTimeout:  getEnvAsDuration("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getEnvAsDuration("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getEnvAsDuration("IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			Name:     getEnv("DB_NAME", "opensame"),
			User:     getEnv("DB_USER", "opensame"),
			Password: getEnv("DB_PASSWORD", "opensame_password"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvAsInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
		},
		RabbitMQ: RabbitMQConfig{
			Host:     getEnv("RABBITMQ_HOST", "localhost"),
			Port:     getEnvAsInt("RABBITMQ_PORT", 5672),
			User:     getEnv("RABBITMQ_USER", "opensame"),
			Password: getEnv("RABBITMQ_PASS", "opensame_password"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 24),
			RefreshHours:     getEnvAsInt("JWT_REFRESH_HOURS", 168), // 7 days
		},
		AI: AIConfig{
			OpenAIKey:      getEnv("OPENAI_API_KEY", ""),
			OpenAIModel:    getEnv("OPENAI_MODEL", "gpt-4"),
			AnthropicKey:   getEnv("ANTHROPIC_API_KEY", ""),
			AnthropicModel: getEnv("ANTHROPIC_MODEL", "claude-3-sonnet-20240229"),
			MaxTokens:      getEnvAsInt("AI_MAX_TOKENS", 4000),
			Temperature:    getEnvAsFloat("AI_TEMPERATURE", 0.7),
		},
		RateLimit: getEnvAsFloat("RATE_LIMIT", 100.0), // requests per second
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsFloat(key string, defaultValue float64) float64 {
	if value := os.Getenv(key); value != "" {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
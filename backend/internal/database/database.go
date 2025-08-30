package database

import (
	"fmt"
	"log"
	"time"

	"github.com/open-same/backend/internal/config"
	"github.com/open-same/backend/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init initializes the database connection
func Init(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get underlying sql.DB
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Database connection established successfully")

	// Auto migrate models
	if err := AutoMigrate(); err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %v", err)
	}

	return DB, nil
}

// AutoMigrate automatically migrates the database schema
func AutoMigrate() error {
	log.Println("Starting database migration...")

	// Enable UUID extension
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %v", err)
	}

	// Enable JSONB extension
	if err := DB.Exec("CREATE EXTENSION IF NOT EXISTS \"pg_trgm\"").Error; err != nil {
		return fmt.Errorf("failed to create pg_trgm extension: %v", err)
	}

	// Migrate models
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Token{},
		&models.Content{},
		&models.ContentVersion{},
		&models.SharedContent{},
		&models.Collaboration{},
	}

	for _, model := range modelsToMigrate {
		if err := DB.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}

	log.Println("Database migration completed successfully")
	return nil
}

// CreateIndexes creates additional database indexes for performance
func CreateIndexes() error {
	log.Println("Creating database indexes...")

	// User indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error; err != nil {
		return fmt.Errorf("failed to create users email index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error; err != nil {
		return fmt.Errorf("failed to create users username index: %v", err)
	}

	// Content indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_content_user_id ON contents(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create content user_id index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_content_type ON contents(type)").Error; err != nil {
		return fmt.Errorf("failed to create content type index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_content_status ON contents(status)").Error; err != nil {
		return fmt.Errorf("failed to create content status index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_content_public ON contents(is_public)").Error; err != nil {
		return fmt.Errorf("failed to create content public index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_content_tags ON contents USING GIN(tags)").Error; err != nil {
		return fmt.Errorf("failed to create content tags index: %v", err)
	}

	// Collaboration indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_collaborations_content_id ON collaborations(content_id)").Error; err != nil {
		return fmt.Errorf("failed to create collaborations content_id index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_collaborations_user_id ON collaborations(user_id)").Error; err != nil {
		return fmt.Errorf("failed to create collaborations user_id index: %v", err)
	}

	// Shared content indexes
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_shared_contents_content_id ON shared_contents(content_id)").Error; err != nil {
		return fmt.Errorf("failed to create shared_contents content_id index: %v", err)
	}
	if err := DB.Exec("CREATE INDEX IF NOT EXISTS idx_shared_contents_shared_with ON shared_contents(shared_with)").Error; err != nil {
		return fmt.Errorf("failed to create shared_contents shared_with index: %v", err)
	}

	log.Println("Database indexes created successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// Transaction executes a function within a database transaction
func Transaction(fn func(tx *gorm.DB) error) error {
	return DB.Transaction(fn)
}
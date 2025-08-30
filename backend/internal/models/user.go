package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID                uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email             string         `json:"email" gorm:"uniqueIndex;not null"`
	Username          string         `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash      string         `json:"-" gorm:"not null"`
	FirstName         string         `json:"first_name"`
	LastName          string         `json:"last_name"`
	Avatar            string         `json:"avatar"`
	Bio               string         `json:"bio"`
	IsVerified        bool           `json:"is_verified" gorm:"default:false"`
	IsActive          bool           `json:"is_active" gorm:"default:true"`
	IsAdmin           bool           `json:"is_admin" gorm:"default:false"`
	LastLoginAt       *time.Time     `json:"last_login_at"`
	EmailVerifiedAt   *time.Time     `json:"email_verified_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relationships
	Contents          []Content      `json:"contents,omitempty" gorm:"foreignKey:UserID"`
	Collaborations    []Collaboration `json:"collaborations,omitempty" gorm:"foreignKey:UserID"`
	SharedContents    []SharedContent `json:"shared_contents,omitempty" gorm:"foreignKey:OwnerID"`
	Tokens            []Token        `json:"tokens,omitempty" gorm:"foreignKey:UserID"`
}

// Token represents user authentication tokens
type Token struct {
	ID           uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID       uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Token        string         `json:"token" gorm:"uniqueIndex;not null"`
	Type         string         `json:"type" gorm:"not null"` // access, refresh, reset
	ExpiresAt    time.Time      `json:"expires_at" gorm:"not null"`
	IsRevoked    bool           `json:"is_revoked" gorm:"default:false"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	
	// Relationships
	User         User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook to set timestamps
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

// BeforeCreate hook for Token
func (t *Token) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// SetPassword hashes and sets the user's password
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword verifies the user's password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// FullName returns the user's full name
func (u *User) FullName() string {
	if u.FirstName != "" && u.LastName != "" {
		return u.FirstName + " " + u.LastName
	}
	if u.FirstName != "" {
		return u.FirstName
	}
	return u.Username
}

// IsTokenExpired checks if a token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// Revoke marks a token as revoked
func (t *Token) Revoke() {
	t.IsRevoked = true
	t.UpdatedAt = time.Now()
}
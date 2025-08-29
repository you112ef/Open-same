package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a platform user
type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	Username     string    `json:"username" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Avatar      string    `json:"avatar"`
	Bio         string    `json:"bio"`
	IsAdmin     bool      `json:"is_admin" gorm:"default:false"`
	IsVerified  bool      `json:"is_verified" gorm:"default:false"`
	IsBanned    bool      `json:"is_banned" gorm:"default:false"`
	LastLogin   time.Time `json:"last_login"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Content        []Content        `json:"content,omitempty" gorm:"foreignKey:CreatorID"`
	Collaborations []Collaboration  `json:"collaborations,omitempty" gorm:"foreignKey:UserID"`
	SharedContent  []SharedContent  `json:"shared_content,omitempty" gorm:"foreignKey:SharedByID"`
}

// Content represents user-created content
type Content struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Type        string    `json:"type" gorm:"not null"` // document, code, diagram, etc.
	Content     string    `json:"content" gorm:"type:text"`
	Metadata    JSON      `json:"metadata" gorm:"type:jsonb"`
	IsPublic   bool      `json:"is_public" gorm:"default:false"`
	IsTemplate  bool      `json:"is_template" gorm:"default:false"`
	Version     int       `json:"version" gorm:"default:1"`
	CreatorID   uuid.UUID `json:"creator_id" gorm:"type:uuid;not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	// Relationships
	Creator       User            `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	Collaborators []Collaboration `json:"collaborators,omitempty" gorm:"foreignKey:ContentID"`
	Shared        []SharedContent  `json:"shared,omitempty" gorm:"foreignKey:ContentID"`
	Versions      []ContentVersion `json:"versions,omitempty" gorm:"foreignKey:ContentID"`
}

// ContentVersion represents version history of content
type ContentVersion struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID uuid.UUID `json:"content_id" gorm:"type:uuid;not null"`
	Version   int       `json:"version" gorm:"not null"`
	Title     string    `json:"title"`
	Content   string    `json:"content" gorm:"type:text"`
	Metadata  JSON      `json:"metadata" gorm:"type:jsonb"`
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Content  Content `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	Creator  User    `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
}

// Collaboration represents user collaboration on content
type Collaboration struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	ContentID uuid.UUID `json:"content_id" gorm:"type:uuid;not null"`
	Role      string    `json:"role" gorm:"not null"` // owner, editor, viewer
	Status    string    `json:"status" gorm:"default:pending"` // pending, accepted, declined
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User    User    `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content Content `json:"content,omitempty" gorm:"foreignKey:ContentID"`
}

// SharedContent represents publicly shared content
type SharedContent struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID uuid.UUID `json:"content_id" gorm:"type:uuid;not null"`
	SharedBy  uuid.UUID `json:"shared_by" gorm:"type:uuid;not null"`
	ShareType string    `json:"share_type" gorm:"not null"` // public, link, embed
	ShareURL  string    `json:"share_url" gorm:"uniqueIndex"`
	ExpiresAt *time.Time `json:"expires_at"`
	ViewCount int        `json:"view_count" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Content  Content `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	SharedBy User    `json:"shared_by_user,omitempty" gorm:"foreignKey:SharedBy"`
}

// JSON is a custom type for JSONB fields
type JSON map[string]interface{}

// BeforeCreate hooks for GORM
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

func (c *Content) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

func (cv *ContentVersion) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == uuid.Nil {
		cv.ID = uuid.New()
	}
	return nil
}

func (col *Collaboration) BeforeCreate(tx *gorm.DB) error {
	if col.ID == uuid.Nil {
		col.ID = uuid.New()
	}
	return nil
}

func (sc *SharedContent) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}
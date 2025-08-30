package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContentType represents the type of content
type ContentType string

const (
	ContentTypeText     ContentType = "text"
	ContentTypeCode     ContentType = "code"
	ContentTypeDiagram  ContentType = "diagram"
	ContentTypeImage    ContentType = "image"
	ContentTypeDocument ContentType = "document"
	ContentTypeTemplate ContentType = "template"
)

// ContentStatus represents the status of content
type ContentStatus string

const (
	ContentStatusDraft     ContentStatus = "draft"
	ContentStatusPublished ContentStatus = "published"
	ContentStatusArchived  ContentStatus = "archived"
	ContentStatusDeleted   ContentStatus = "deleted"
)

// Content represents user-generated content
type Content struct {
	ID              uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID          uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Title           string         `json:"title" gorm:"not null"`
	Description     string         `json:"description"`
	Content         string         `json:"content" gorm:"type:text"`
	Type            ContentType    `json:"type" gorm:"not null;default:'text'"`
	Status          ContentStatus  `json:"status" gorm:"not null;default:'draft'"`
	IsPublic        bool           `json:"is_public" gorm:"default:false"`
	IsTemplate      bool           `json:"is_template" gorm:"default:false"`
	Tags            []string       `json:"tags" gorm:"type:text[]"`
	Metadata        JSON           `json:"metadata" gorm:"type:jsonb"`
	AIGenerated     bool           `json:"ai_generated" gorm:"default:false"`
	AIModel         string         `json:"ai_model"`
	AIPrompt        string         `json:"ai_prompt"`
	Version         int            `json:"version" gorm:"default:1"`
	ParentID        *uuid.UUID     `json:"parent_id" gorm:"type:uuid"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	
	// Relationships
	User            User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Parent          *Content       `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Versions        []ContentVersion `json:"versions,omitempty" gorm:"foreignKey:ContentID"`
	Collaborations  []Collaboration `json:"collaborations,omitempty" gorm:"foreignKey:ContentID"`
	SharedContents  []SharedContent `json:"shared_contents,omitempty" gorm:"foreignKey:ContentID"`
}

// ContentVersion represents a version of content
type ContentVersion struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID   uuid.UUID      `json:"content_id" gorm:"type:uuid;not null"`
	Version     int            `json:"version" gorm:"not null"`
	Content     string         `json:"content" gorm:"type:text;not null"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Tags        []string       `json:"tags" gorm:"type:text[]"`
	Metadata    JSON           `json:"metadata" gorm:"type:jsonb"`
	CreatedBy   uuid.UUID      `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt   time.Time      `json:"created_at"`
	
	// Relationships
	Content      Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	User        User           `json:"user,omitempty" gorm:"foreignKey:CreatedBy"`
}

// SharedContent represents content shared with other users
type SharedContent struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID   uuid.UUID      `json:"content_id" gorm:"type:uuid;not null"`
	OwnerID     uuid.UUID      `json:"owner_id" gorm:"type:uuid;not null"`
	SharedWith  uuid.UUID      `json:"shared_with" gorm:"type:uuid;not null"`
	Permission  string         `json:"permission" gorm:"not null;default:'read'"` // read, write, admin
	ExpiresAt   *time.Time     `json:"expires_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	
	// Relationships
	Content     Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	Owner      User           `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	SharedUser User           `json:"shared_user,omitempty" gorm:"foreignKey:SharedWith"`
}

// Collaboration represents user collaboration on content
type Collaboration struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ContentID   uuid.UUID      `json:"content_id" gorm:"type:uuid;not null"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	Role        string         `json:"role" gorm:"not null;default:'editor'"` // viewer, editor, admin
	JoinedAt    time.Time      `json:"joined_at"`
	LastActive  *time.Time     `json:"last_active"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	
	// Relationships
	Content     Content        `json:"content,omitempty" gorm:"foreignKey:ContentID"`
	User       User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// JSON is a custom type for JSONB fields
type JSON map[string]interface{}

// BeforeCreate hooks
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

func (sc *SharedContent) BeforeCreate(tx *gorm.DB) error {
	if sc.ID == uuid.Nil {
		sc.ID = uuid.New()
	}
	return nil
}

func (col *Collaboration) BeforeCreate(tx *gorm.DB) error {
	if col.ID == uuid.Nil {
		col.ID = uuid.New()
	}
	return nil
}

// GetFullContent returns the full content with metadata
func (c *Content) GetFullContent() string {
	return c.Content
}

// IsCollaborator checks if a user is a collaborator
func (c *Content) IsCollaborator(userID uuid.UUID) bool {
	for _, col := range c.Collaborations {
		if col.UserID == userID && col.IsActive {
			return true
		}
	}
	return false
}

// CanEdit checks if a user can edit the content
func (c *Content) CanEdit(userID uuid.UUID) bool {
	if c.UserID == userID {
		return true
	}
	
	for _, col := range c.Collaborations {
		if col.UserID == userID && col.IsActive && (col.Role == "editor" || col.Role == "admin") {
			return true
		}
	}
	return false
}

// CanAdmin checks if a user can admin the content
func (c *Content) CanAdmin(userID uuid.UUID) bool {
	if c.UserID == userID {
		return true
	}
	
	for _, col := range c.Collaborations {
		if col.UserID == userID && col.IsActive && col.Role == "admin" {
			return true
		}
	}
	return false
}
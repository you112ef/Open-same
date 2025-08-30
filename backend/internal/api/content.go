package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/open-same/backend/internal/database"
	"github.com/open-same/backend/internal/middleware"
	"github.com/open-same/backend/internal/models"
)

// CreateContentRequest represents content creation request
type CreateContentRequest struct {
	Title       string                `json:"title" binding:"required,min=1,max=200"`
	Description string                `json:"description"`
	Content     string                `json:"content"`
	Type        models.ContentType    `json:"type" binding:"required"`
	IsPublic    bool                  `json:"is_public"`
	IsTemplate  bool                  `json:"is_template"`
	Tags        []string              `json:"tags"`
	Metadata    map[string]interface{} `json:"metadata"`
	ParentID    *string               `json:"parent_id"`
}

// UpdateContentRequest represents content update request
type UpdateContentRequest struct {
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	Content     *string                `json:"content"`
	Type        *models.ContentType    `json:"type"`
	Status      *models.ContentStatus  `json:"status"`
	IsPublic    *bool                  `json:"is_public"`
	IsTemplate  *bool                  `json:"is_template"`
	Tags        *[]string              `json:"tags"`
	Metadata    *map[string]interface{} `json:"metadata"`
}

// ContentListResponse represents paginated content list response
type ContentListResponse struct {
	Contents    []models.Content `json:"contents"`
	Total       int64            `json:"total"`
	Page        int              `json:"page"`
	PerPage     int              `json:"per_page"`
	TotalPages  int              `json:"total_pages"`
	HasNext     bool             `json:"has_next"`
	HasPrevious bool             `json:"has_previous"`
}

// CreateContent handles content creation
func CreateContent(c *gin.Context) {
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// Get user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "User context not found",
			"code":    "MISSING_USER_CONTEXT",
			"message": "Internal server error",
		})
		return
	}

	// Parse parent ID if provided
	var parentID *uuid.UUID
	if req.ParentID != nil {
		parsedID, err := uuid.Parse(*req.ParentID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid parent ID",
				"code":    "INVALID_PARENT_ID",
				"message": "Parent ID must be a valid UUID",
			})
			return
		}
		parentID = &parsedID
	}

	// Create content
	content := models.Content{
		UserID:      user.ID,
		Title:       req.Title,
		Description: req.Description,
		Content:     req.Content,
		Type:        req.Type,
		Status:      models.ContentStatusDraft,
		IsPublic:    req.IsPublic,
		IsTemplate:  req.IsTemplate,
		Tags:        req.Tags,
		Metadata:    models.JSON(req.Metadata),
		ParentID:    parentID,
		Version:     1,
	}

	// Save content to database
	if err := database.GetDB().Create(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create content",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while creating content",
		})
		return
	}

	// Create initial version
	version := models.ContentVersion{
		ContentID:   content.ID,
		Version:     1,
		Content:     content.Content,
		Title:       content.Title,
		Description: content.Description,
		Tags:        content.Tags,
		Metadata:    content.Metadata,
		CreatedBy:   user.ID,
	}

	if err := database.GetDB().Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create content version",
			"code":    "VERSION_CREATION_ERROR",
			"message": "Content created but version tracking failed",
		})
		return
	}

	// Load relationships
	database.GetDB().Preload("User").First(&content, content.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message": "Content created successfully",
		"data":    content,
	})
}

// GetContent handles content retrieval
func GetContent(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Content ID required",
			"code":    "MISSING_CONTENT_ID",
			"message": "Content ID is required",
		})
		return
	}

	// Parse content ID
	id, err := uuid.Parse(contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid content ID",
			"code":    "INVALID_CONTENT_ID",
			"message": "Content ID must be a valid UUID",
		})
		return
	}

	// Get content with relationships
	var content models.Content
	if err := database.GetDB().Preload("User").Preload("Versions").Preload("Collaborations.User").First(&content, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Content not found",
			"code":    "CONTENT_NOT_FOUND",
			"message": "The requested content was not found",
		})
		return
	}

	// Check if user can access this content
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		// Public content can be accessed without authentication
		if !content.IsPublic {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"code":    "ACCESS_DENIED",
				"message": "You don't have permission to access this content",
			})
			return
		}
	} else {
		// Check if user owns the content or is a collaborator
		if content.UserID != user.ID && !content.IsCollaborator(user.ID) && !content.IsPublic {
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Access denied",
				"code":    "ACCESS_DENIED",
				"message": "You don't have permission to access this content",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Content retrieved successfully",
		"data":    content,
	})
}

// GetUserContent handles user content list retrieval
func GetUserContent(c *gin.Context) {
	// Get user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "User context not found",
			"code":    "MISSING_USER_CONTEXT",
			"message": "Internal server error",
		})
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	contentType := c.Query("type")
	status := c.Query("status")
	search := c.Query("search")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Build query
	query := database.GetDB().Model(&models.Content{}).Where("user_id = ?", user.ID)

	// Apply filters
	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Calculate pagination
	offset := (page - 1) * perPage
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	// Get content with pagination
	var contents []models.Content
	if err := query.Preload("User").Offset(offset).Limit(perPage).Order("updated_at DESC").Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve content",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while retrieving content",
		})
		return
	}

	response := ContentListResponse{
		Contents:    contents,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Content retrieved successfully",
		"data":    response,
	})
}

// UpdateContent handles content updates
func UpdateContent(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Content ID required",
			"code":    "MISSING_CONTENT_ID",
			"message": "Content ID is required",
		})
		return
	}

	// Parse content ID
	id, err := uuid.Parse(contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid content ID",
			"code":    "INVALID_CONTENT_ID",
			"message": "Content ID must be a valid UUID",
		})
		return
	}

	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request data",
			"code":    "INVALID_REQUEST",
			"message": err.Error(),
		})
		return
	}

	// Get user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "User context not found",
			"code":    "MISSING_USER_CONTEXT",
			"message": "Internal server error",
		})
		return
	}

	// Get content
	var content models.Content
	if err := database.GetDB().First(&content, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Content not found",
			"code":    "CONTENT_NOT_FOUND",
			"message": "The requested content was not found",
		})
		return
	}

	// Check if user can edit this content
	if !content.CanEdit(user.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Edit permission denied",
			"code":    "EDIT_PERMISSION_DENIED",
			"message": "You don't have permission to edit this content",
		})
		return
	}

	// Create new version if content changed
	contentChanged := false
	if req.Content != nil && *req.Content != content.Content {
		contentChanged = true
	}

	// Update fields
	if req.Title != nil {
		content.Title = *req.Title
		contentChanged = true
	}
	if req.Description != nil {
		content.Description = *req.Description
		contentChanged = true
	}
	if req.Content != nil {
		content.Content = *req.Content
		contentChanged = true
	}
	if req.Type != nil {
		content.Type = *req.Type
		contentChanged = true
	}
	if req.Status != nil {
		content.Status = *req.Status
		contentChanged = true
	}
	if req.IsPublic != nil {
		content.IsPublic = *req.IsPublic
		contentChanged = true
	}
	if req.IsTemplate != nil {
		content.IsTemplate = *req.IsTemplate
		contentChanged = true
	}
	if req.Tags != nil {
		content.Tags = *req.Tags
		contentChanged = true
	}
	if req.Metadata != nil {
		content.Metadata = models.JSON(*req.Metadata)
		contentChanged = true
	}

	// Update timestamp
	content.UpdatedAt = time.Now()

	// Save content
	if err := database.GetDB().Save(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update content",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while updating content",
		})
		return
	}

	// Create new version if content changed
	if contentChanged {
		content.Version++
		version := models.ContentVersion{
			ContentID:   content.ID,
			Version:     content.Version,
			Content:     content.Content,
			Title:       content.Title,
			Description: content.Description,
			Tags:        content.Tags,
			Metadata:    content.Metadata,
			CreatedBy:   user.ID,
		}

		if err := database.GetDB().Create(&version).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Failed to create content version",
				"code":    "VERSION_CREATION_ERROR",
				"message": "Content updated but version tracking failed",
			})
			return
		}
	}

	// Load relationships
	database.GetDB().Preload("User").First(&content, content.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Content updated successfully",
		"data":    content,
	})
}

// DeleteContent handles content deletion
func DeleteContent(c *gin.Context) {
	contentID := c.Param("id")
	if contentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Content ID required",
			"code":    "MISSING_CONTENT_ID",
			"message": "Content ID is required",
		})
		return
	}

	// Parse content ID
	id, err := uuid.Parse(contentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid content ID",
			"code":    "INVALID_CONTENT_ID",
			"message": "Content ID must be a valid UUID",
		})
		return
	}

	// Get user from context
	user, exists := middleware.GetUserFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "User context not found",
			"code":    "MISSING_USER_CONTEXT",
			"message": "Internal server error",
		})
		return
	}

	// Get content
	var content models.Content
	if err := database.GetDB().First(&content, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Content not found",
			"code":    "CONTENT_NOT_FOUND",
			"message": "The requested content was not found",
		})
		return
	}

	// Check if user can delete this content
	if !content.CanAdmin(user.ID) {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Delete permission denied",
			"code":    "DELETE_PERMISSION_DENIED",
			"message": "You don't have permission to delete this content",
		})
		return
	}

	// Soft delete content
	if err := database.GetDB().Delete(&content).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete content",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while deleting content",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Content deleted successfully",
	})
}

// GetPublicContent handles public content retrieval
func GetPublicContent(c *gin.Context) {
	// Parse query parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))
	contentType := c.Query("type")
	search := c.Query("search")

	// Validate pagination
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	// Build query for public content
	query := database.GetDB().Model(&models.Content{}).Where("is_public = ? AND status = ?", true, models.ContentStatusPublished)

	// Apply filters
	if contentType != "" {
		query = query.Where("type = ?", contentType)
	}
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Calculate pagination
	offset := (page - 1) * perPage
	totalPages := int((total + int64(perPage) - 1) / int64(perPage))

	// Get content with pagination
	var contents []models.Content
	if err := query.Preload("User").Offset(offset).Limit(perPage).Order("created_at DESC").Find(&contents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to retrieve content",
			"code":    "DATABASE_ERROR",
			"message": "An error occurred while retrieving content",
		})
		return
	}

	response := ContentListResponse{
		Contents:    contents,
		Total:       total,
		Page:        page,
		PerPage:     perPage,
		TotalPages:  totalPages,
		HasNext:     page < totalPages,
		HasPrevious: page > 1,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Public content retrieved successfully",
		"data":    response,
	})
}
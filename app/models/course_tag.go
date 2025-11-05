package models

import (
	"time"

	"gorm.io/gorm"
)

// CourseTag represents a courseTag entity
type CourseTag struct {
	Id        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name      string         `json:"name"`
}

// TableName returns the table name for the CourseTag model
func (m *CourseTag) TableName() string {
	return "course_tags"
}

// GetId returns the Id of the model
func (m *CourseTag) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseTag) GetModelName() string {
	return "course_tag"
}

// CreateCourseTagRequest represents the request payload for creating a CourseTag
type CreateCourseTagRequest struct {
	Name string `json:"name"`
}

// UpdateCourseTagRequest represents the request payload for updating a CourseTag
type UpdateCourseTagRequest struct {
	Name string `json:"name,omitempty"`
}

// CourseTagResponse represents the API response for CourseTag
type CourseTagResponse struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
}

// CourseTagModelResponse represents a simplified response when this model is part of other entities
type CourseTagModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// CourseTagSelectOption represents a simplified response for select boxes and dropdowns
type CourseTagSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // From Name field
}

// CourseTagListResponse represents the response for list operations (optimized for performance)
type CourseTagListResponse struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Name      string         `json:"name"`
}

// ToResponse converts the model to an API response
func (m *CourseTag) ToResponse() *CourseTagResponse {
	if m == nil {
		return nil
	}
	response := &CourseTagResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		Name:      m.Name,
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseTag) ToModelResponse() *CourseTagModelResponse {
	if m == nil {
		return nil
	}
	return &CourseTagModelResponse{
		Id:   m.Id,
		Name: m.Name,
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseTag) ToSelectOption() *CourseTagSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Name

	return &CourseTagSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseTag) ToListResponse() *CourseTagListResponse {
	if m == nil {
		return nil
	}
	return &CourseTagListResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		Name:      m.Name,
	}
}

// Preload preloads all the model's relationships
func (m *CourseTag) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}

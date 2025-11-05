package models

import (
	"base/core/types"
	"time"

	"gorm.io/gorm"
)

// CourseResource represents a courseResource entity
type CourseResource struct {
	Id         uint           `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	FileUrl    string         `json:"file_url"`
	Title      string         `json:"title"`
	UploadedAt types.DateTime `json:"uploaded_at"`
	CourseId   uint           `json:"course_id,omitempty"`
	Course     *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
}

// TableName returns the table name for the CourseResource model
func (m *CourseResource) TableName() string {
	return "course_resources"
}

// GetId returns the Id of the model
func (m *CourseResource) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseResource) GetModelName() string {
	return "course_resource"
}

// CreateCourseResourceRequest represents the request payload for creating a CourseResource
type CreateCourseResourceRequest struct {
	CourseId   uint           `json:"course_id,omitempty"`
	FileUrl    string         `json:"file_url"`
	Title      string         `json:"title"`
	UploadedAt types.DateTime `json:"uploaded_at" swaggertype:"string"`
}

// UpdateCourseResourceRequest represents the request payload for updating a CourseResource
type UpdateCourseResourceRequest struct {
	CourseId   uint           `json:"course_id,omitempty"`
	FileUrl    string         `json:"file_url,omitempty"`
	Title      string         `json:"title,omitempty"`
	UploadedAt types.DateTime `json:"uploaded_at,omitempty" swaggertype:"string"`
}

// CourseResourceResponse represents the API response for CourseResource
type CourseResourceResponse struct {
	Id         uint                 `json:"id"`
	CreatedAt  time.Time            `json:"created_at"`
	UpdatedAt  time.Time            `json:"updated_at"`
	DeletedAt  gorm.DeletedAt       `json:"deleted_at"`
	FileUrl    string               `json:"file_url"`
	Title      string               `json:"title"`
	UploadedAt types.DateTime       `json:"uploaded_at"`
	Course     *CourseModelResponse `json:"course,omitempty"`
}

// CourseResourceModelResponse represents a simplified response when this model is part of other entities
type CourseResourceModelResponse struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

// CourseResourceSelectOption represents a simplified response for select boxes and dropdowns
type CourseResourceSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // From Title field
}

// CourseResourceListResponse represents the response for list operations (optimized for performance)
type CourseResourceListResponse struct {
	Id         uint           `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
	FileUrl    string         `json:"file_url"`
	Title      string         `json:"title"`
	UploadedAt types.DateTime `json:"uploaded_at"`
}

// ToResponse converts the model to an API response
func (m *CourseResource) ToResponse() *CourseResourceResponse {
	if m == nil {
		return nil
	}
	response := &CourseResourceResponse{
		Id:         m.Id,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
		FileUrl:    m.FileUrl,
		Title:      m.Title,
		UploadedAt: m.UploadedAt,
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseResource) ToModelResponse() *CourseResourceModelResponse {
	if m == nil {
		return nil
	}
	return &CourseResourceModelResponse{
		Id:    m.Id,
		Title: m.Title,
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseResource) ToSelectOption() *CourseResourceSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Title

	return &CourseResourceSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseResource) ToListResponse() *CourseResourceListResponse {
	if m == nil {
		return nil
	}
	return &CourseResourceListResponse{
		Id:         m.Id,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
		FileUrl:    m.FileUrl,
		Title:      m.Title,
		UploadedAt: m.UploadedAt,
	}
}

// Preload preloads all the model's relationships
func (m *CourseResource) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Course")
	return query
}

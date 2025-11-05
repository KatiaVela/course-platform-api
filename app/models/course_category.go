package models

import (
	"time"

	"gorm.io/gorm"
)

// CourseCategory represents a courseCategory entity
type CourseCategory struct {
	Id          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Courses     []Course       `json:"courses,omitempty" gorm:"foreignKey:CategoryId"`
}

// TableName returns the table name for the CourseCategory model
func (m *CourseCategory) TableName() string {
	return "course_categories"
}

// GetId returns the Id of the model
func (m *CourseCategory) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseCategory) GetModelName() string {
	return "course_category"
}

// CreateCourseCategoryRequest represents the request payload for creating a CourseCategory
type CreateCourseCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// UpdateCourseCategoryRequest represents the request payload for updating a CourseCategory
type UpdateCourseCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
}

// CourseCategoryResponse represents the API response for CourseCategory
type CourseCategoryResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	Course      []*Course      `json:"course,omitempty"`
}

// CourseCategoryModelResponse represents a simplified response when this model is part of other entities
type CourseCategoryModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// CourseCategorySelectOption represents a simplified response for select boxes and dropdowns
type CourseCategorySelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // From Name field
}

// CourseCategoryListResponse represents the response for list operations (optimized for performance)
type CourseCategoryListResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
}

// ToResponse converts the model to an API response
func (m *CourseCategory) ToResponse() *CourseCategoryResponse {
	if m == nil {
		return nil
	}
	response := &CourseCategoryResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseCategory) ToModelResponse() *CourseCategoryModelResponse {
	if m == nil {
		return nil
	}
	return &CourseCategoryModelResponse{
		Id:   m.Id,
		Name: m.Name,
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseCategory) ToSelectOption() *CourseCategorySelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Name

	return &CourseCategorySelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseCategory) ToListResponse() *CourseCategoryListResponse {
	if m == nil {
		return nil
	}
	return &CourseCategoryListResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		Name:        m.Name,
		Slug:        m.Slug,
		Description: m.Description,
	}
}

// Preload preloads all the model's relationships
func (m *CourseCategory) Preload(db *gorm.DB) *gorm.DB {
	query := db
	return query
}

package models

import (
	"base/core/app/profile"
	"time"

	"gorm.io/gorm"
)

// Course represents a course entity
type Course struct {
	Id           uint            `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at" gorm:"index"`
	Title        string          `json:"title"`
	Slug         string          `json:"slug"`
	Description  string          `json:"description"`
	Price        int             `json:"price"`
	Level        string          `json:"level"`
	Language     string          `json:"language"`
	ThumbnailUrl string          `json:"thumbnail_url"`
	Status       string          `json:"status"`
	Duration     int             `json:"duration"`
	InstructorId uint            `json:"instructor_id,omitempty"`
	CategoryId   *uint           `json:"category_id,omitempty" gorm:"index"`
	Instructor   *profile.User   `json:"instructor,omitempty" gorm:"foreignKey:InstructorId"`
	Category     *CourseCategory `json:"category,omitempty" gorm:"foreignKey:CategoryId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// TableName returns the table name for the Course model
func (m *Course) TableName() string {
	return "courses"
}

// GetId returns the Id of the model
func (m *Course) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *Course) GetModelName() string {
	return "course"
}

// CreateCourseRequest represents the request payload for creating a Course
type CreateCourseRequest struct {
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	InstructorId uint   `json:"instructor_id,omitempty"`
	CategoryId   *uint  `json:"category_id,omitempty"`
	Price        int    `json:"price"`
	Level        string `json:"level"`
	Language     string `json:"language"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Status       string `json:"status"`
	Duration     int    `json:"duration"`
}

// UpdateCourseRequest represents the request payload for updating a Course
type UpdateCourseRequest struct {
	Title        string `json:"title,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Description  string `json:"description,omitempty"`
	InstructorId uint   `json:"instructor_id,omitempty"`
	CategoryId   *uint  `json:"category_id,omitempty"`
	Price        int    `json:"price,omitempty"`
	Level        string `json:"level,omitempty"`
	Language     string `json:"language,omitempty"`
	ThumbnailUrl string `json:"thumbnail_url,omitempty"`
	Status       string `json:"status,omitempty"`
	Duration     int    `json:"duration,omitempty"`
}

// CourseResponse represents the API response for Course
type CourseResponse struct {
	Id           uint                         `json:"id"`
	CreatedAt    time.Time                    `json:"created_at"`
	UpdatedAt    time.Time                    `json:"updated_at"`
	DeletedAt    gorm.DeletedAt               `json:"deleted_at"`
	Title        string                       `json:"title"`
	Slug         string                       `json:"slug"`
	Description  string                       `json:"description"`
	Price        int                          `json:"price"`
	Level        string                       `json:"level"`
	Language     string                       `json:"language"`
	ThumbnailUrl string                       `json:"thumbnail_url"`
	Status       string                       `json:"status"`
	Duration     int                          `json:"duration"`
	Instructor   *profile.UserModelResponse   `json:"instructor,omitempty"`
	Category     *CourseCategoryModelResponse `json:"category,omitempty"`
}

// CourseModelResponse represents a simplified response when this model is part of other entities
type CourseModelResponse struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

// CourseSelectOption represents a simplified response for select boxes and dropdowns
type CourseSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // From Title field
}

// CourseListResponse represents the response for list operations (optimized for performance)
type CourseListResponse struct {
	Id           uint           `json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at"`
	Title        string         `json:"title"`
	Slug         string         `json:"slug"`
	Description  string         `json:"description"`
	Price        int            `json:"price"`
	Level        string         `json:"level"`
	Language     string         `json:"language"`
	ThumbnailUrl string         `json:"thumbnail_url"`
	Status       string         `json:"status"`
	Duration     int            `json:"duration"`
}

// ToResponse converts the model to an API response
func (m *Course) ToResponse() *CourseResponse {
	if m == nil {
		return nil
	}
	response := &CourseResponse{
		Id:           m.Id,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
		Title:        m.Title,
		Slug:         m.Slug,
		Description:  m.Description,
		Price:        m.Price,
		Level:        m.Level,
		Language:     m.Language,
		ThumbnailUrl: m.ThumbnailUrl,
		Status:       m.Status,
		Duration:     m.Duration,
	}
	if m.InstructorId != 0 {
		response.Instructor = m.Instructor.ToModelResponse()
	}
	if m.CategoryId != nil {
		response.Category = m.Category.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *Course) ToModelResponse() *CourseModelResponse {
	if m == nil {
		return nil
	}
	return &CourseModelResponse{
		Id:    m.Id,
		Title: m.Title,
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *Course) ToSelectOption() *CourseSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Title

	return &CourseSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *Course) ToListResponse() *CourseListResponse {
	if m == nil {
		return nil
	}
	return &CourseListResponse{
		Id:           m.Id,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
		Title:        m.Title,
		Slug:         m.Slug,
		Description:  m.Description,
		Price:        m.Price,
		Level:        m.Level,
		Language:     m.Language,
		ThumbnailUrl: m.ThumbnailUrl,
		Status:       m.Status,
		Duration:     m.Duration,
	}
}

// Preload preloads all the model's relationships
func (m *Course) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Instructor")
	query = query.Preload("Category")
	return query
}

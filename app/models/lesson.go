package models

import (
	"time"

	"gorm.io/gorm"
)

// Lesson represents a lesson entity
type Lesson struct {
	Id          uint           `json:"id" gorm:"primarykey"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	VideoUrl    string         `json:"video_url"`
	Duration    int            `json:"duration"`
	OrderNumber int            `json:"order_number"`
	CourseId    uint           `json:"course_id,omitempty"`
	Course      *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
}

// TableName returns the table name for the Lesson model
func (m *Lesson) TableName() string {
	return "lessons"
}

// GetId returns the Id of the model
func (m *Lesson) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *Lesson) GetModelName() string {
	return "lesson"
}

// CreateLessonRequest represents the request payload for creating a Lesson
type CreateLessonRequest struct {
	Title       string `json:"title"`
	CourseId    uint   `json:"course_id,omitempty"`
	Content     string `json:"content"`
	VideoUrl    string `json:"video_url"`
	Duration    int    `json:"duration"`
	OrderNumber int    `json:"order_number"`
}

// UpdateLessonRequest represents the request payload for updating a Lesson
type UpdateLessonRequest struct {
	Title       string `json:"title,omitempty"`
	CourseId    uint   `json:"course_id,omitempty"`
	Content     string `json:"content,omitempty"`
	VideoUrl    string `json:"video_url,omitempty"`
	Duration    int    `json:"duration,omitempty"`
	OrderNumber int    `json:"order_number,omitempty"`
}

// LessonResponse represents the API response for Lesson
type LessonResponse struct {
	Id          uint                 `json:"id"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   gorm.DeletedAt       `json:"deleted_at"`
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	VideoUrl    string               `json:"video_url"`
	Duration    int                  `json:"duration"`
	OrderNumber int                  `json:"order_number"`
	Course      *CourseModelResponse `json:"course,omitempty"`
}

// LessonModelResponse represents a simplified response when this model is part of other entities
type LessonModelResponse struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

// LessonSelectOption represents a simplified response for select boxes and dropdowns
type LessonSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // From Title field
}

// LessonListResponse represents the response for list operations (optimized for performance)
type LessonListResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	VideoUrl    string         `json:"video_url"`
	Duration    int            `json:"duration"`
	OrderNumber int            `json:"order_number"`
}

// ToResponse converts the model to an API response
func (m *Lesson) ToResponse() *LessonResponse {
	if m == nil {
		return nil
	}
	response := &LessonResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		Title:       m.Title,
		Content:     m.Content,
		VideoUrl:    m.VideoUrl,
		Duration:    m.Duration,
		OrderNumber: m.OrderNumber,
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *Lesson) ToModelResponse() *LessonModelResponse {
	if m == nil {
		return nil
	}
	return &LessonModelResponse{
		Id:    m.Id,
		Title: m.Title,
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *Lesson) ToSelectOption() *LessonSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Title

	return &LessonSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *Lesson) ToListResponse() *LessonListResponse {
	if m == nil {
		return nil
	}
	return &LessonListResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		Title:       m.Title,
		Content:     m.Content,
		VideoUrl:    m.VideoUrl,
		Duration:    m.Duration,
		OrderNumber: m.OrderNumber,
	}
}

// Preload preloads all the model's relationships
func (m *Lesson) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Course")
	return query
}

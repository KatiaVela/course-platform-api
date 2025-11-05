package models

import (
	"base/core/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CourseProgressLog represents a courseProgressLog entity
type CourseProgressLog struct {
	Id           uint           `json:"id" gorm:"primarykey"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CompletedAt  types.DateTime `json:"completed_at"`
	EnrollmentId uint           `json:"enrollment_id,omitempty"`
	LessonId     uint           `json:"lesson_id,omitempty"`
	Enrollment   *Enrollment    `json:"enrollment,omitempty" gorm:"foreignKey:EnrollmentId"`
	Lesson       *Lesson        `json:"lesson,omitempty" gorm:"foreignKey:LessonId"`
}

// TableName returns the table name for the CourseProgressLog model
func (m *CourseProgressLog) TableName() string {
	return "course_progress_logs"
}

// GetId returns the Id of the model
func (m *CourseProgressLog) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseProgressLog) GetModelName() string {
	return "course_progress_log"
}

// CreateCourseProgressLogRequest represents the request payload for creating a CourseProgressLog
type CreateCourseProgressLogRequest struct {
	EnrollmentId uint           `json:"enrollment_id,omitempty"`
	LessonId     uint           `json:"lesson_id,omitempty"`
	CompletedAt  types.DateTime `json:"completed_at" swaggertype:"string"`
}

// UpdateCourseProgressLogRequest represents the request payload for updating a CourseProgressLog
type UpdateCourseProgressLogRequest struct {
	EnrollmentId uint           `json:"enrollment_id,omitempty"`
	LessonId     uint           `json:"lesson_id,omitempty"`
	CompletedAt  types.DateTime `json:"completed_at,omitempty" swaggertype:"string"`
}

// CourseProgressLogResponse represents the API response for CourseProgressLog
type CourseProgressLogResponse struct {
	Id          uint                     `json:"id"`
	CreatedAt   time.Time                `json:"created_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	DeletedAt   gorm.DeletedAt           `json:"deleted_at"`
	CompletedAt types.DateTime           `json:"completed_at"`
	Enrollment  *EnrollmentModelResponse `json:"enrollment,omitempty"`
	Lesson      *LessonModelResponse     `json:"lesson,omitempty"`
}

// CourseProgressLogModelResponse represents a simplified response when this model is part of other entities
type CourseProgressLogModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseProgressLogSelectOption represents a simplified response for select boxes and dropdowns
type CourseProgressLogSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseProgressLogListResponse represents the response for list operations (optimized for performance)
type CourseProgressLogListResponse struct {
	Id          uint           `json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
	CompletedAt types.DateTime `json:"completed_at"`
}

// ToResponse converts the model to an API response
func (m *CourseProgressLog) ToResponse() *CourseProgressLogResponse {
	if m == nil {
		return nil
	}
	response := &CourseProgressLogResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		CompletedAt: m.CompletedAt,
	}
	if m.EnrollmentId != 0 {
		response.Enrollment = m.Enrollment.ToModelResponse()
	}
	if m.LessonId != 0 {
		response.Lesson = m.Lesson.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseProgressLog) ToModelResponse() *CourseProgressLogModelResponse {
	if m == nil {
		return nil
	}
	return &CourseProgressLogModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("CourseProgressLog #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseProgressLog) ToSelectOption() *CourseProgressLogSelectOption {
	if m == nil {
		return nil
	}
	displayName := fmt.Sprintf("CourseProgressLog #%d", m.Id) // Fallback to ID-based display

	return &CourseProgressLogSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseProgressLog) ToListResponse() *CourseProgressLogListResponse {
	if m == nil {
		return nil
	}
	return &CourseProgressLogListResponse{
		Id:          m.Id,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   m.DeletedAt,
		CompletedAt: m.CompletedAt,
	}
}

// Preload preloads all the model's relationships
func (m *CourseProgressLog) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Enrollment")
	query = query.Preload("Lesson")
	return query
}

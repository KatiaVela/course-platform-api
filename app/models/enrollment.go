package models

import (
	"base/core/app/profile"
	"base/core/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Enrollment represents a enrollment entity
type Enrollment struct {
	Id         uint           `json:"id" gorm:"primarykey"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	EnrolledAt types.DateTime `json:"enrolled_at"`
	Progress   int            `json:"progress"`
	Completed  bool           `json:"completed"`
	StudentId  uint           `json:"student_id,omitempty"`
	CourseId   uint           `json:"course_id,omitempty"`
	Student    *profile.User  `json:"student,omitempty" gorm:"foreignKey:StudentId"`
	Course     *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
}

// TableName returns the table name for the Enrollment model
func (m *Enrollment) TableName() string {
	return "enrollments"
}

// GetId returns the Id of the model
func (m *Enrollment) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *Enrollment) GetModelName() string {
	return "enrollment"
}

// CreateEnrollmentRequest represents the request payload for creating a Enrollment
type CreateEnrollmentRequest struct {
	StudentId  uint           `json:"student_id,omitempty"`
	CourseId   uint           `json:"course_id,omitempty"`
	EnrolledAt types.DateTime `json:"enrolled_at" swaggertype:"string"`
	Progress   int            `json:"progress"`
	Completed  bool           `json:"completed"`
}

// UpdateEnrollmentRequest represents the request payload for updating a Enrollment
type UpdateEnrollmentRequest struct {
	StudentId  uint           `json:"student_id,omitempty"`
	CourseId   uint           `json:"course_id,omitempty"`
	EnrolledAt types.DateTime `json:"enrolled_at,omitempty" swaggertype:"string"`
	Progress   int            `json:"progress,omitempty"`
	Completed  *bool          `json:"completed,omitempty"`
}

// EnrollmentResponse represents the API response for Enrollment
type EnrollmentResponse struct {
	Id         uint                       `json:"id"`
	CreatedAt  time.Time                  `json:"created_at"`
	UpdatedAt  time.Time                  `json:"updated_at"`
	DeletedAt  gorm.DeletedAt             `json:"deleted_at"`
	EnrolledAt types.DateTime             `json:"enrolled_at"`
	Progress   int                        `json:"progress"`
	Completed  bool                       `json:"completed"`
	Student    *profile.UserModelResponse `json:"student,omitempty"`
	Course     *CourseModelResponse       `json:"course,omitempty"`
}

// EnrollmentModelResponse represents a simplified response when this model is part of other entities
type EnrollmentModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// EnrollmentSelectOption represents a simplified response for select boxes and dropdowns
type EnrollmentSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// EnrollmentListResponse represents the response for list operations (optimized for performance)
type EnrollmentListResponse struct {
	Id         uint           `json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
	EnrolledAt types.DateTime `json:"enrolled_at"`
	Progress   int            `json:"progress"`
	Completed  bool           `json:"completed"`
}

// ToResponse converts the model to an API response
func (m *Enrollment) ToResponse() *EnrollmentResponse {
	if m == nil {
		return nil
	}
	response := &EnrollmentResponse{
		Id:         m.Id,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
		EnrolledAt: m.EnrolledAt,
		Progress:   m.Progress,
		Completed:  m.Completed,
	}
	if m.StudentId != 0 {
		response.Student = m.Student.ToModelResponse()
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *Enrollment) ToModelResponse() *EnrollmentModelResponse {
	if m == nil {
		return nil
	}
	return &EnrollmentModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("Enrollment #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *Enrollment) ToSelectOption() *EnrollmentSelectOption {
	if m == nil {
		return nil
	}
	displayName := fmt.Sprintf("Enrollment #%d", m.Id) // Fallback to ID-based display

	return &EnrollmentSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *Enrollment) ToListResponse() *EnrollmentListResponse {
	if m == nil {
		return nil
	}
	return &EnrollmentListResponse{
		Id:         m.Id,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		DeletedAt:  m.DeletedAt,
		EnrolledAt: m.EnrolledAt,
		Progress:   m.Progress,
		Completed:  m.Completed,
	}
}

// Preload preloads all the model's relationships
func (m *Enrollment) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Student")
	query = query.Preload("Course")
	return query
}

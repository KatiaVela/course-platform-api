package models

import (
	"base/core/app/profile"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Review represents a review entity
type Review struct {
	Id        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Rating    int            `json:"rating"`
	Comment   string         `json:"comment"`
	CourseId  uint           `json:"course_id,omitempty"`
	StudentId uint           `json:"student_id,omitempty"`
	Course    *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
	Student   *profile.User  `json:"student,omitempty" gorm:"foreignKey:StudentId"`
}

// TableName returns the table name for the Review model
func (m *Review) TableName() string {
	return "reviews"
}

// GetId returns the Id of the model
func (m *Review) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *Review) GetModelName() string {
	return "review"
}

// CreateReviewRequest represents the request payload for creating a Review
type CreateReviewRequest struct {
	CourseId  uint   `json:"course_id,omitempty"`
	StudentId uint   `json:"student_id,omitempty"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

// UpdateReviewRequest represents the request payload for updating a Review
type UpdateReviewRequest struct {
	CourseId  uint   `json:"course_id,omitempty"`
	StudentId uint   `json:"student_id,omitempty"`
	Rating    int    `json:"rating,omitempty"`
	Comment   string `json:"comment,omitempty"`
}

// ReviewResponse represents the API response for Review
type ReviewResponse struct {
	Id        uint                       `json:"id"`
	CreatedAt time.Time                  `json:"created_at"`
	UpdatedAt time.Time                  `json:"updated_at"`
	DeletedAt gorm.DeletedAt             `json:"deleted_at"`
	Rating    int                        `json:"rating"`
	Comment   string                     `json:"comment"`
	Course    *CourseModelResponse       `json:"course,omitempty"`
	Student   *profile.UserModelResponse `json:"student,omitempty"`
}

// ReviewModelResponse represents a simplified response when this model is part of other entities
type ReviewModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// ReviewSelectOption represents a simplified response for select boxes and dropdowns
type ReviewSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// ReviewListResponse represents the response for list operations (optimized for performance)
type ReviewListResponse struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Rating    int            `json:"rating"`
	Comment   string         `json:"comment"`
}

// ToResponse converts the model to an API response
func (m *Review) ToResponse() *ReviewResponse {
	if m == nil {
		return nil
	}
	response := &ReviewResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		Rating:    m.Rating,
		Comment:   m.Comment,
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}
	if m.StudentId != 0 {
		response.Student = m.Student.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *Review) ToModelResponse() *ReviewModelResponse {
	if m == nil {
		return nil
	}
	return &ReviewModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("Review #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *Review) ToSelectOption() *ReviewSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.Comment // Using first string field as display name

	return &ReviewSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *Review) ToListResponse() *ReviewListResponse {
	if m == nil {
		return nil
	}
	return &ReviewListResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
		Rating:    m.Rating,
		Comment:   m.Comment,
	}
}

// Preload preloads all the model's relationships
func (m *Review) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Course")
	query = query.Preload("Student")
	return query
}

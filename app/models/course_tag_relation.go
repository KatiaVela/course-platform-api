package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CourseTagRelation represents a courseTagRelation entity
type CourseTagRelation struct {
	Id        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CourseId  uint           `json:"course_id,omitempty"`
	TagId     uint           `json:"tag_id,omitempty"`
	Course    *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
	Tag       *CourseTag     `json:"tag,omitempty" gorm:"foreignKey:TagId"`
}

// TableName returns the table name for the CourseTagRelation model
func (m *CourseTagRelation) TableName() string {
	return "course_tag_relations"
}

// GetId returns the Id of the model
func (m *CourseTagRelation) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseTagRelation) GetModelName() string {
	return "course_tag_relation"
}

// CreateCourseTagRelationRequest represents the request payload for creating a CourseTagRelation
type CreateCourseTagRelationRequest struct {
	CourseId uint `json:"course_id,omitempty"`
	TagId    uint `json:"tag_id,omitempty"`
}

// UpdateCourseTagRelationRequest represents the request payload for updating a CourseTagRelation
type UpdateCourseTagRelationRequest struct {
	CourseId uint `json:"course_id,omitempty"`
	TagId    uint `json:"tag_id,omitempty"`
}

// CourseTagRelationResponse represents the API response for CourseTagRelation
type CourseTagRelationResponse struct {
	Id        uint                    `json:"id"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
	DeletedAt gorm.DeletedAt          `json:"deleted_at"`
	Course    *CourseModelResponse    `json:"course,omitempty"`
	Tag       *CourseTagModelResponse `json:"tag,omitempty"`
}

// CourseTagRelationModelResponse represents a simplified response when this model is part of other entities
type CourseTagRelationModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseTagRelationSelectOption represents a simplified response for select boxes and dropdowns
type CourseTagRelationSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseTagRelationListResponse represents the response for list operations (optimized for performance)
type CourseTagRelationListResponse struct {
	Id        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

// ToResponse converts the model to an API response
func (m *CourseTagRelation) ToResponse() *CourseTagRelationResponse {
	if m == nil {
		return nil
	}
	response := &CourseTagRelationResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}
	if m.TagId != 0 {
		response.Tag = m.Tag.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseTagRelation) ToModelResponse() *CourseTagRelationModelResponse {
	if m == nil {
		return nil
	}
	return &CourseTagRelationModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("CourseTagRelation #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseTagRelation) ToSelectOption() *CourseTagRelationSelectOption {
	if m == nil {
		return nil
	}
	displayName := fmt.Sprintf("CourseTagRelation #%d", m.Id) // Fallback to ID-based display

	return &CourseTagRelationSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseTagRelation) ToListResponse() *CourseTagRelationListResponse {
	if m == nil {
		return nil
	}
	return &CourseTagRelationListResponse{
		Id:        m.Id,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt,
	}
}

// Preload preloads all the model's relationships
func (m *CourseTagRelation) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Course")
	query = query.Preload("Tag")
	return query
}

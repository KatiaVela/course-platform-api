package models

import (
	"base/core/types"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// CourseCertificate represents a courseCertificate entity
type CourseCertificate struct {
	Id             uint           `json:"id" gorm:"primarykey"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	CertificateUrl string         `json:"certificate_url"`
	IssuedAt       types.DateTime `json:"issued_at"`
	EnrollmentId   uint           `json:"enrollment_id,omitempty"`
	Enrollment     *Enrollment    `json:"enrollment,omitempty" gorm:"foreignKey:EnrollmentId"`
}

// TableName returns the table name for the CourseCertificate model
func (m *CourseCertificate) TableName() string {
	return "course_certificates"
}

// GetId returns the Id of the model
func (m *CourseCertificate) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *CourseCertificate) GetModelName() string {
	return "course_certificate"
}

// CreateCourseCertificateRequest represents the request payload for creating a CourseCertificate
type CreateCourseCertificateRequest struct {
	EnrollmentId   uint           `json:"enrollment_id,omitempty"`
	CertificateUrl string         `json:"certificate_url"`
	IssuedAt       types.DateTime `json:"issued_at" swaggertype:"string"`
}

// UpdateCourseCertificateRequest represents the request payload for updating a CourseCertificate
type UpdateCourseCertificateRequest struct {
	EnrollmentId   uint           `json:"enrollment_id,omitempty"`
	CertificateUrl string         `json:"certificate_url,omitempty"`
	IssuedAt       types.DateTime `json:"issued_at,omitempty" swaggertype:"string"`
}

// CourseCertificateResponse represents the API response for CourseCertificate
type CourseCertificateResponse struct {
	Id             uint                     `json:"id"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	DeletedAt      gorm.DeletedAt           `json:"deleted_at"`
	CertificateUrl string                   `json:"certificate_url"`
	IssuedAt       types.DateTime           `json:"issued_at"`
	Enrollment     *EnrollmentModelResponse `json:"enrollment,omitempty"`
}

// CourseCertificateModelResponse represents a simplified response when this model is part of other entities
type CourseCertificateModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseCertificateSelectOption represents a simplified response for select boxes and dropdowns
type CourseCertificateSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// CourseCertificateListResponse represents the response for list operations (optimized for performance)
type CourseCertificateListResponse struct {
	Id             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
	CertificateUrl string         `json:"certificate_url"`
	IssuedAt       types.DateTime `json:"issued_at"`
}

// ToResponse converts the model to an API response
func (m *CourseCertificate) ToResponse() *CourseCertificateResponse {
	if m == nil {
		return nil
	}
	response := &CourseCertificateResponse{
		Id:             m.Id,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
		CertificateUrl: m.CertificateUrl,
		IssuedAt:       m.IssuedAt,
	}
	if m.EnrollmentId != 0 {
		response.Enrollment = m.Enrollment.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *CourseCertificate) ToModelResponse() *CourseCertificateModelResponse {
	if m == nil {
		return nil
	}
	return &CourseCertificateModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("CourseCertificate #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *CourseCertificate) ToSelectOption() *CourseCertificateSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.CertificateUrl // Using first string field as display name

	return &CourseCertificateSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *CourseCertificate) ToListResponse() *CourseCertificateListResponse {
	if m == nil {
		return nil
	}
	return &CourseCertificateListResponse{
		Id:             m.Id,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      m.DeletedAt,
		CertificateUrl: m.CertificateUrl,
		IssuedAt:       m.IssuedAt,
	}
}

// Preload preloads all the model's relationships
func (m *CourseCertificate) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("Enrollment")
	return query
}

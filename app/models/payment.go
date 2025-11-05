package models

import (
	"base/core/app/profile"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type PaymentMethod string

const (
	PaymentMethodCreditCard   PaymentMethod = "credit_card"
	PaymentMethodPaypal       PaymentMethod = "paypal"
	PaymentMethodBankTransfer PaymentMethod = "bank_transfer"
)

type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
)

// Payment represents a payment entity
type Payment struct {
	Id            uint           `json:"id" gorm:"primarykey"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Amount        int            `json:"amount"`
	PaymentMethod PaymentMethod  `json:"payment_method" gorm:"type:ENUM('credit_card', 'paypal', 'bank_transfer')"`
	PaymentStatus PaymentStatus  `json:"payment_status" gorm:"type:ENUM('pending', 'completed', 'failed')"`
	TransactionId string         `json:"transaction_id"`
	UserId        uint           `json:"user_id,omitempty"`
	CourseId      uint           `json:"course_id,omitempty"`
	User          *profile.User  `json:"user,omitempty" gorm:"foreignKey:UserId"`
	Course        *Course        `json:"course,omitempty" gorm:"foreignKey:CourseId"`
}

// TableName returns the table name for the Payment model
func (m *Payment) TableName() string {
	return "payments"
}

// GetId returns the Id of the model
func (m *Payment) GetId() uint {
	return m.Id
}

// GetModelName returns the model name
func (m *Payment) GetModelName() string {
	return "payment"
}

// CreatePaymentRequest represents the request payload for creating a Payment
type CreatePaymentRequest struct {
	UserId        uint   `json:"user_id,omitempty"`
	CourseId      uint   `json:"course_id,omitempty"`
	Amount        int           `json:"amount"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	PaymentStatus PaymentStatus `json:"payment_status"`
	TransactionId string        `json:"transaction_id"`
}

// UpdatePaymentRequest represents the request payload for updating a Payment
type UpdatePaymentRequest struct {
	UserId        uint   `json:"user_id,omitempty"`
	CourseId      uint   `json:"course_id,omitempty"`
	Amount        int           `json:"amount,omitempty"`
	PaymentMethod PaymentMethod `json:"payment_method,omitempty"`
	PaymentStatus PaymentStatus `json:"payment_status,omitempty"`
	TransactionId string        `json:"transaction_id,omitempty"`
}

// PaymentResponse represents the API response for Payment
type PaymentResponse struct {
	Id            uint                       `json:"id"`
	CreatedAt     time.Time                  `json:"created_at"`
	UpdatedAt     time.Time                  `json:"updated_at"`
	DeletedAt     gorm.DeletedAt             `json:"deleted_at"`
	Amount        int                        `json:"amount"`
	PaymentMethod PaymentMethod             `json:"payment_method"`
	PaymentStatus PaymentStatus             `json:"payment_status"`
	TransactionId string                     `json:"transaction_id"`
	User          *profile.UserModelResponse `json:"user,omitempty"`
	Course        *CourseModelResponse       `json:"course,omitempty"`
}

// PaymentModelResponse represents a simplified response when this model is part of other entities
type PaymentModelResponse struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// PaymentSelectOption represents a simplified response for select boxes and dropdowns
type PaymentSelectOption struct {
	Id   uint   `json:"id"`
	Name string `json:"name"` // Display name
}

// PaymentListResponse represents the response for list operations (optimized for performance)
type PaymentListResponse struct {
	Id            uint           `json:"id"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
	Amount        int            `json:"amount"`
	PaymentMethod PaymentMethod  `json:"payment_method"`
	PaymentStatus PaymentStatus  `json:"payment_status"`
	TransactionId string         `json:"transaction_id"`
}

// ToResponse converts the model to an API response
func (m *Payment) ToResponse() *PaymentResponse {
	if m == nil {
		return nil
	}
	response := &PaymentResponse{
		Id:            m.Id,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
		Amount:        m.Amount,
		PaymentMethod: m.PaymentMethod,
		PaymentStatus: m.PaymentStatus,
		TransactionId: m.TransactionId,
	}
	if m.UserId != 0 {
		response.User = m.User.ToModelResponse()
	}
	if m.CourseId != 0 {
		response.Course = m.Course.ToModelResponse()
	}

	return response
}

// ToModelResponse converts the model to a simplified response for when it's part of other entities
func (m *Payment) ToModelResponse() *PaymentModelResponse {
	if m == nil {
		return nil
	}
	return &PaymentModelResponse{
		Id:   m.Id,
		Name: fmt.Sprintf("Payment #%d", m.Id), // Fallback to ID-based display
	}
}

// ToSelectOption converts the model to a select option for dropdowns
func (m *Payment) ToSelectOption() *PaymentSelectOption {
	if m == nil {
		return nil
	}
	displayName := m.TransactionId // Using first string field as display name

	return &PaymentSelectOption{
		Id:   m.Id,
		Name: displayName,
	}
}

// ToListResponse converts the model to a list response (without preloaded relationships for fast listing)
func (m *Payment) ToListResponse() *PaymentListResponse {
	if m == nil {
		return nil
	}
	return &PaymentListResponse{
		Id:            m.Id,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     m.DeletedAt,
		Amount:        m.Amount,
		PaymentMethod: m.PaymentMethod,
		PaymentStatus: m.PaymentStatus,
		TransactionId: m.TransactionId,
	}
}

// Preload preloads all the model's relationships
func (m *Payment) Preload(db *gorm.DB) *gorm.DB {
	query := db
	query = query.Preload("User")
	query = query.Preload("Course")
	return query
}

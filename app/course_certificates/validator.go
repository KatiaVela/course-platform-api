package course_certificates

import (
	"base/app/models"
	"base/core/validator"
)

// Global validator instance using Base core validator wrapper
var validate = validator.New()

// ValidateCourseCertificateCreateRequest validates the create request
func ValidateCourseCertificateCreateRequest(req *models.CreateCourseCertificateRequest) error {
	if req == nil {
		return validator.ValidationErrors{
			{
				Field:   "request",
				Tag:     "required",
				Value:   "nil",
				Message: "request cannot be nil",
			},
		}
	}

	// Use Base core validator
	return validate.Validate(req)
}

// ValidateCourseCertificateUpdateRequest validates the update request
func ValidateCourseCertificateUpdateRequest(req *models.UpdateCourseCertificateRequest, id uint) error {
	if req == nil {
		return validator.ValidationErrors{
			{
				Field:   "request",
				Tag:     "required",
				Value:   "nil",
				Message: "request cannot be nil",
			},
		}
	}

	if id == 0 {
		return validator.ValidationErrors{
			{
				Field:   "id",
				Tag:     "required",
				Value:   "0",
				Message: "id cannot be zero",
			},
		}
	}

	// Skip validation for update requests - all fields are optional
	return nil
}

// ValidateCourseCertificateDeleteRequest validates the delete request
func ValidateCourseCertificateDeleteRequest(id uint) error {
	return ValidateID(id)
}

// ValidateID validates if the ID is valid
func ValidateID(id uint) error {
	if id == 0 {
		return validator.ValidationErrors{
			{
				Field:   "id",
				Tag:     "required",
				Value:   "0",
				Message: "id cannot be zero",
			},
		}
	}
	return nil
}

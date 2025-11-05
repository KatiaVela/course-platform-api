package enrollments

import (
	"math"

	"base/app/models"
	"base/core/emitter"
	"base/core/logger"
	"base/core/storage"
	"base/core/types"

	"gorm.io/gorm"
)

const (
	CreateEnrollmentEvent = "enrollments.create"
	UpdateEnrollmentEvent = "enrollments.update"
	DeleteEnrollmentEvent = "enrollments.delete"
)

type EnrollmentService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewEnrollmentService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *EnrollmentService {
	return &EnrollmentService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *EnrollmentService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for Enrollment
	validSortFields := map[string]string{
		"id":          "id",
		"created_at":  "created_at",
		"updated_at":  "updated_at",
		"enrolled_at": "enrolled_at",
		"progress":    "progress",
		"completed":   "completed",
	}

	// Default sorting - if sort_order exists, always use it for custom ordering
	defaultSortBy := "id"
	defaultSortOrder := "desc"

	// Determine sort field
	sortField := defaultSortBy
	if sortBy != nil && *sortBy != "" {
		if field, exists := validSortFields[*sortBy]; exists {
			sortField = field
		}
	}

	// Determine sort direction (order parameter)
	sortDirection := defaultSortOrder
	if sortOrder != nil && (*sortOrder == "asc" || *sortOrder == "desc") {
		sortDirection = *sortOrder
	}

	// Apply sorting
	query.Order(sortField + " " + sortDirection)
}

func (s *EnrollmentService) Create(req *models.CreateEnrollmentRequest) (*models.Enrollment, error) {
	item := &models.Enrollment{
		StudentId:  req.StudentId,
		CourseId:   req.CourseId,
		EnrolledAt: req.EnrolledAt,
		Progress:   req.Progress,
		Completed:  req.Completed,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create enrollment", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateEnrollmentEvent, item)

	return s.GetById(item.Id)
}

func (s *EnrollmentService) Update(id uint, req *models.UpdateEnrollmentRequest) (*models.Enrollment, error) {
	item := &models.Enrollment{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find enrollment for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateEnrollmentUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.StudentId != 0 {
		item.StudentId = req.StudentId
	}
	// For foreign key relationships
	if req.CourseId != 0 {
		item.CourseId = req.CourseId
	}
	// For custom DateTime fields
	if !req.EnrolledAt.IsZero() {
		item.EnrolledAt = req.EnrolledAt
	}
	// For non-pointer integer fields
	if req.Progress != 0 {
		item.Progress = req.Progress
	}
	// For boolean fields, check if it's included in the request (pointer would be non-nil)
	if req.Completed != nil {
		item.Completed = *req.Completed
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update enrollment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated enrollment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateEnrollmentEvent, result)

	return result, nil
}

func (s *EnrollmentService) Delete(id uint) error {
	item := &models.Enrollment{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find enrollment for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete enrollment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteEnrollmentEvent, item)

	return nil
}

func (s *EnrollmentService) GetById(id uint) (*models.Enrollment, error) {
	item := &models.Enrollment{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get enrollment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *EnrollmentService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.Enrollment
	var total int64

	query := s.DB.Model(&models.Enrollment{})
	// Set default values if nil
	defaultPage := 1
	defaultLimit := 10
	if page == nil {
		page = &defaultPage
	}
	if limit == nil {
		limit = &defaultLimit
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		s.Logger.Error("failed to count enrollments",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Apply pagination if provided
	if page != nil && limit != nil {
		offset := (*page - 1) * *limit
		query = query.Offset(offset).Limit(*limit)
	}

	// Apply sorting
	s.applySorting(query, sortBy, sortOrder)

	// Don't preload relationships for list response (faster)
	// query = (&models.Enrollment{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get enrollments",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.EnrollmentListResponse, len(items))
	for i, item := range items {
		responses[i] = item.ToListResponse()
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(*limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	return &types.PaginatedResponse{
		Data: responses,
		Pagination: types.Pagination{
			Total:      int(total),
			Page:       *page,
			PageSize:   *limit,
			TotalPages: totalPages,
		},
	}, nil
}

// GetAllForSelect gets all items for select box/dropdown options (simplified response)
func (s *EnrollmentService) GetAllForSelect() ([]*models.Enrollment, error) {
	var items []*models.Enrollment

	query := s.DB.Model(&models.Enrollment{})

	// Only select the necessary fields for select options
	query = query.Select("id") // Only ID if no name/title field found

	// Order by name/title for better UX
	query = query.Order("id ASC")

	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("Failed to fetch items for select", logger.String("error", err.Error()))
		return nil, err
	}

	return items, nil
}

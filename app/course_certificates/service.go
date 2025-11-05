package course_certificates

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
	CreateCourseCertificateEvent = "coursecertificates.create"
	UpdateCourseCertificateEvent = "coursecertificates.update"
	DeleteCourseCertificateEvent = "coursecertificates.delete"
)

type CourseCertificateService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewCourseCertificateService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *CourseCertificateService {
	return &CourseCertificateService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *CourseCertificateService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for CourseCertificate
	validSortFields := map[string]string{
		"id":              "id",
		"created_at":      "created_at",
		"updated_at":      "updated_at",
		"certificate_url": "certificate_url",
		"issued_at":       "issued_at",
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

func (s *CourseCertificateService) Create(req *models.CreateCourseCertificateRequest) (*models.CourseCertificate, error) {
	item := &models.CourseCertificate{
		EnrollmentId:   req.EnrollmentId,
		CertificateUrl: req.CertificateUrl,
		IssuedAt:       req.IssuedAt,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create coursecertificate", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateCourseCertificateEvent, item)

	return s.GetById(item.Id)
}

func (s *CourseCertificateService) Update(id uint, req *models.UpdateCourseCertificateRequest) (*models.CourseCertificate, error) {
	item := &models.CourseCertificate{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursecertificate for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateCourseCertificateUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.EnrollmentId != 0 {
		item.EnrollmentId = req.EnrollmentId
	}
	// For non-pointer string fields
	if req.CertificateUrl != "" {
		item.CertificateUrl = req.CertificateUrl
	}
	// For custom DateTime fields
	if !req.IssuedAt.IsZero() {
		item.IssuedAt = req.IssuedAt
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update coursecertificate",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated coursecertificate",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateCourseCertificateEvent, result)

	return result, nil
}

func (s *CourseCertificateService) Delete(id uint) error {
	item := &models.CourseCertificate{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursecertificate for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete coursecertificate",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteCourseCertificateEvent, item)

	return nil
}

func (s *CourseCertificateService) GetById(id uint) (*models.CourseCertificate, error) {
	item := &models.CourseCertificate{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get coursecertificate",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *CourseCertificateService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.CourseCertificate
	var total int64

	query := s.DB.Model(&models.CourseCertificate{})
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
		s.Logger.Error("failed to count coursecertificates",
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
	// query = (&models.CourseCertificate{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get coursecertificates",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.CourseCertificateListResponse, len(items))
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
func (s *CourseCertificateService) GetAllForSelect() ([]*models.CourseCertificate, error) {
	var items []*models.CourseCertificate

	query := s.DB.Model(&models.CourseCertificate{})

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

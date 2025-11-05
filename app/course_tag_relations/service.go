package course_tag_relations

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
	CreateCourseTagRelationEvent = "coursetagrelations.create"
	UpdateCourseTagRelationEvent = "coursetagrelations.update"
	DeleteCourseTagRelationEvent = "coursetagrelations.delete"
)

type CourseTagRelationService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewCourseTagRelationService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *CourseTagRelationService {
	return &CourseTagRelationService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *CourseTagRelationService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for CourseTagRelation
	validSortFields := map[string]string{
		"id":         "id",
		"created_at": "created_at",
		"updated_at": "updated_at",
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

func (s *CourseTagRelationService) Create(req *models.CreateCourseTagRelationRequest) (*models.CourseTagRelation, error) {
	item := &models.CourseTagRelation{
		CourseId: req.CourseId,
		TagId:    req.TagId,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create coursetagrelation", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateCourseTagRelationEvent, item)

	return s.GetById(item.Id)
}

func (s *CourseTagRelationService) Update(id uint, req *models.UpdateCourseTagRelationRequest) (*models.CourseTagRelation, error) {
	item := &models.CourseTagRelation{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursetagrelation for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateCourseTagRelationUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.CourseId != 0 {
		item.CourseId = req.CourseId
	}
	// For foreign key relationships
	if req.TagId != 0 {
		item.TagId = req.TagId
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update coursetagrelation",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated coursetagrelation",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateCourseTagRelationEvent, result)

	return result, nil
}

func (s *CourseTagRelationService) Delete(id uint) error {
	item := &models.CourseTagRelation{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursetagrelation for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete coursetagrelation",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteCourseTagRelationEvent, item)

	return nil
}

func (s *CourseTagRelationService) GetById(id uint) (*models.CourseTagRelation, error) {
	item := &models.CourseTagRelation{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get coursetagrelation",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *CourseTagRelationService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.CourseTagRelation
	var total int64

	query := s.DB.Model(&models.CourseTagRelation{})
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
		s.Logger.Error("failed to count coursetagrelations",
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
	// query = (&models.CourseTagRelation{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get coursetagrelations",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.CourseTagRelationListResponse, len(items))
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
func (s *CourseTagRelationService) GetAllForSelect() ([]*models.CourseTagRelation, error) {
	var items []*models.CourseTagRelation

	query := s.DB.Model(&models.CourseTagRelation{})

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

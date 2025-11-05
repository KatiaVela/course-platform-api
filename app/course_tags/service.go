package course_tags

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
	CreateCourseTagEvent = "coursetags.create"
	UpdateCourseTagEvent = "coursetags.update"
	DeleteCourseTagEvent = "coursetags.delete"
)

type CourseTagService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewCourseTagService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *CourseTagService {
	return &CourseTagService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *CourseTagService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for CourseTag
	validSortFields := map[string]string{
		"id":         "id",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"name":       "name",
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

func (s *CourseTagService) Create(req *models.CreateCourseTagRequest) (*models.CourseTag, error) {
	item := &models.CourseTag{
		Name: req.Name,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create coursetag", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateCourseTagEvent, item)

	return s.GetById(item.Id)
}

func (s *CourseTagService) Update(id uint, req *models.UpdateCourseTagRequest) (*models.CourseTag, error) {
	item := &models.CourseTag{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursetag for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateCourseTagUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For non-pointer string fields
	if req.Name != "" {
		item.Name = req.Name
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update coursetag",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated coursetag",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateCourseTagEvent, result)

	return result, nil
}

func (s *CourseTagService) Delete(id uint) error {
	item := &models.CourseTag{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find coursetag for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete coursetag",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteCourseTagEvent, item)

	return nil
}

func (s *CourseTagService) GetById(id uint) (*models.CourseTag, error) {
	item := &models.CourseTag{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get coursetag",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *CourseTagService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.CourseTag
	var total int64

	query := s.DB.Model(&models.CourseTag{})
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
		s.Logger.Error("failed to count coursetags",
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
	// query = (&models.CourseTag{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get coursetags",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.CourseTagListResponse, len(items))
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
func (s *CourseTagService) GetAllForSelect() ([]*models.CourseTag, error) {
	var items []*models.CourseTag

	query := s.DB.Model(&models.CourseTag{})

	// Only select the necessary fields for select options
	query = query.Select("id, name")

	// Order by name/title for better UX
	query = query.Order("name ASC")

	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("Failed to fetch items for select", logger.String("error", err.Error()))
		return nil, err
	}

	return items, nil
}

package courses

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
	CreateCourseEvent = "courses.create"
	UpdateCourseEvent = "courses.update"
	DeleteCourseEvent = "courses.delete"
)

type CourseService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewCourseService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *CourseService {
	return &CourseService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *CourseService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for Course
	validSortFields := map[string]string{
		"id":            "id",
		"created_at":    "created_at",
		"updated_at":    "updated_at",
		"title":         "title",
		"slug":          "slug",
		"description":   "description",
		"price":         "price",
		"level":         "level",
		"language":      "language",
		"thumbnail_url": "thumbnail_url",
		"status":        "status",
		"duration":      "duration",
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

func (s *CourseService) Create(req *models.CreateCourseRequest) (*models.Course, error) {
	item := &models.Course{
		Title:        req.Title,
		Slug:         req.Slug,
		Description:  req.Description,
		InstructorId: req.InstructorId,
		CategoryId:   req.CategoryId,
		Price:        req.Price,
		Level:        req.Level,
		Language:     req.Language,
		ThumbnailUrl: req.ThumbnailUrl,
		Status:       req.Status,
		Duration:     req.Duration,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create course", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateCourseEvent, item)

	return s.GetById(item.Id)
}

func (s *CourseService) Update(id uint, req *models.UpdateCourseRequest) (*models.Course, error) {
	item := &models.Course{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find course for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateCourseUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For non-pointer string fields
	if req.Title != "" {
		item.Title = req.Title
	}
	// For non-pointer string fields
	if req.Slug != "" {
		item.Slug = req.Slug
	}
	// For non-pointer string fields
	if req.Description != "" {
		item.Description = req.Description
	}
	// For foreign key relationships
	if req.InstructorId != 0 {
		item.InstructorId = req.InstructorId
	}
	// For foreign key relationships
	if req.CategoryId != nil {
		item.CategoryId = req.CategoryId
	}
	// For non-pointer integer fields
	if req.Price != 0 {
		item.Price = req.Price
	}
	// For non-pointer string fields
	if req.Level != "" {
		item.Level = req.Level
	}
	// For non-pointer string fields
	if req.Language != "" {
		item.Language = req.Language
	}
	// For non-pointer string fields
	if req.ThumbnailUrl != "" {
		item.ThumbnailUrl = req.ThumbnailUrl
	}
	// For non-pointer string fields
	if req.Status != "" {
		item.Status = req.Status
	}
	// For non-pointer integer fields
	if req.Duration != 0 {
		item.Duration = req.Duration
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update course",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated course",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateCourseEvent, result)

	return result, nil
}

func (s *CourseService) Delete(id uint) error {
	item := &models.Course{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find course for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete course",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteCourseEvent, item)

	return nil
}

func (s *CourseService) GetById(id uint) (*models.Course, error) {
	item := &models.Course{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get course",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *CourseService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.Course
	var total int64

	query := s.DB.Model(&models.Course{})
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
		s.Logger.Error("failed to count courses",
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
	// query = (&models.Course{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get courses",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.CourseListResponse, len(items))
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
func (s *CourseService) GetAllForSelect() ([]*models.Course, error) {
	var items []*models.Course

	query := s.DB.Model(&models.Course{})

	// Only select the necessary fields for select options
	query = query.Select("id, title")

	// Order by name/title for better UX
	query = query.Order("title ASC")

	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("Failed to fetch items for select", logger.String("error", err.Error()))
		return nil, err
	}

	return items, nil
}

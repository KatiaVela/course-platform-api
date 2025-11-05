package lessons

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
	CreateLessonEvent = "lessons.create"
	UpdateLessonEvent = "lessons.update"
	DeleteLessonEvent = "lessons.delete"
)

type LessonService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewLessonService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *LessonService {
	return &LessonService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *LessonService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for Lesson
	validSortFields := map[string]string{
		"id":           "id",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"title":        "title",
		"content":      "content",
		"video_url":    "video_url",
		"duration":     "duration",
		"order_number": "order_number",
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

func (s *LessonService) Create(req *models.CreateLessonRequest) (*models.Lesson, error) {
	item := &models.Lesson{
		Title:       req.Title,
		CourseId:    req.CourseId,
		Content:     req.Content,
		VideoUrl:    req.VideoUrl,
		Duration:    req.Duration,
		OrderNumber: req.OrderNumber,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create lesson", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateLessonEvent, item)

	return s.GetById(item.Id)
}

func (s *LessonService) Update(id uint, req *models.UpdateLessonRequest) (*models.Lesson, error) {
	item := &models.Lesson{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find lesson for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateLessonUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For non-pointer string fields
	if req.Title != "" {
		item.Title = req.Title
	}
	// For foreign key relationships
	if req.CourseId != 0 {
		item.CourseId = req.CourseId
	}
	// For non-pointer string fields
	if req.Content != "" {
		item.Content = req.Content
	}
	// For non-pointer string fields
	if req.VideoUrl != "" {
		item.VideoUrl = req.VideoUrl
	}
	// For non-pointer integer fields
	if req.Duration != 0 {
		item.Duration = req.Duration
	}
	// For non-pointer integer fields
	if req.OrderNumber != 0 {
		item.OrderNumber = req.OrderNumber
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update lesson",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated lesson",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateLessonEvent, result)

	return result, nil
}

func (s *LessonService) Delete(id uint) error {
	item := &models.Lesson{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find lesson for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete lesson",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteLessonEvent, item)

	return nil
}

func (s *LessonService) GetById(id uint) (*models.Lesson, error) {
	item := &models.Lesson{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get lesson",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *LessonService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.Lesson
	var total int64

	query := s.DB.Model(&models.Lesson{})
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
		s.Logger.Error("failed to count lessons",
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
	// query = (&models.Lesson{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get lessons",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.LessonListResponse, len(items))
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
func (s *LessonService) GetAllForSelect() ([]*models.Lesson, error) {
	var items []*models.Lesson

	query := s.DB.Model(&models.Lesson{})

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

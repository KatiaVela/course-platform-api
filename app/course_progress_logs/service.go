package course_progress_logs

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
	CreateCourseProgressLogEvent = "courseprogresslogs.create"
	UpdateCourseProgressLogEvent = "courseprogresslogs.update"
	DeleteCourseProgressLogEvent = "courseprogresslogs.delete"
)

type CourseProgressLogService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewCourseProgressLogService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *CourseProgressLogService {
	return &CourseProgressLogService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *CourseProgressLogService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for CourseProgressLog
	validSortFields := map[string]string{
		"id":           "id",
		"created_at":   "created_at",
		"updated_at":   "updated_at",
		"completed_at": "completed_at",
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

func (s *CourseProgressLogService) Create(req *models.CreateCourseProgressLogRequest) (*models.CourseProgressLog, error) {
	item := &models.CourseProgressLog{
		EnrollmentId: req.EnrollmentId,
		LessonId:     req.LessonId,
		CompletedAt:  req.CompletedAt,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create courseprogresslog", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateCourseProgressLogEvent, item)

	return s.GetById(item.Id)
}

func (s *CourseProgressLogService) Update(id uint, req *models.UpdateCourseProgressLogRequest) (*models.CourseProgressLog, error) {
	item := &models.CourseProgressLog{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find courseprogresslog for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateCourseProgressLogUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.EnrollmentId != 0 {
		item.EnrollmentId = req.EnrollmentId
	}
	// For foreign key relationships
	if req.LessonId != 0 {
		item.LessonId = req.LessonId
	}
	// For custom DateTime fields
	if !req.CompletedAt.IsZero() {
		item.CompletedAt = req.CompletedAt
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update courseprogresslog",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated courseprogresslog",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateCourseProgressLogEvent, result)

	return result, nil
}

func (s *CourseProgressLogService) Delete(id uint) error {
	item := &models.CourseProgressLog{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find courseprogresslog for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete courseprogresslog",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteCourseProgressLogEvent, item)

	return nil
}

func (s *CourseProgressLogService) GetById(id uint) (*models.CourseProgressLog, error) {
	item := &models.CourseProgressLog{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get courseprogresslog",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *CourseProgressLogService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.CourseProgressLog
	var total int64

	query := s.DB.Model(&models.CourseProgressLog{})
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
		s.Logger.Error("failed to count courseprogresslogs",
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
	// query = (&models.CourseProgressLog{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get courseprogresslogs",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.CourseProgressLogListResponse, len(items))
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
func (s *CourseProgressLogService) GetAllForSelect() ([]*models.CourseProgressLog, error) {
	var items []*models.CourseProgressLog

	query := s.DB.Model(&models.CourseProgressLog{})

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

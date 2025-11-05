package reviews

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
	CreateReviewEvent = "reviews.create"
	UpdateReviewEvent = "reviews.update"
	DeleteReviewEvent = "reviews.delete"
)

type ReviewService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewReviewService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *ReviewService {
	return &ReviewService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *ReviewService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for Review
	validSortFields := map[string]string{
		"id":         "id",
		"created_at": "created_at",
		"updated_at": "updated_at",
		"rating":     "rating",
		"comment":    "comment",
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

func (s *ReviewService) Create(req *models.CreateReviewRequest) (*models.Review, error) {
	item := &models.Review{
		CourseId:  req.CourseId,
		StudentId: req.StudentId,
		Rating:    req.Rating,
		Comment:   req.Comment,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create review", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreateReviewEvent, item)

	return s.GetById(item.Id)
}

func (s *ReviewService) Update(id uint, req *models.UpdateReviewRequest) (*models.Review, error) {
	item := &models.Review{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find review for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidateReviewUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.CourseId != 0 {
		item.CourseId = req.CourseId
	}
	// For foreign key relationships
	if req.StudentId != 0 {
		item.StudentId = req.StudentId
	}
	// For non-pointer integer fields
	if req.Rating != 0 {
		item.Rating = req.Rating
	}
	// For non-pointer string fields
	if req.Comment != "" {
		item.Comment = req.Comment
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update review",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated review",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdateReviewEvent, result)

	return result, nil
}

func (s *ReviewService) Delete(id uint) error {
	item := &models.Review{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find review for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete review",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeleteReviewEvent, item)

	return nil
}

func (s *ReviewService) GetById(id uint) (*models.Review, error) {
	item := &models.Review{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get review",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *ReviewService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.Review
	var total int64

	query := s.DB.Model(&models.Review{})
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
		s.Logger.Error("failed to count reviews",
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
	// query = (&models.Review{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get reviews",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.ReviewListResponse, len(items))
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
func (s *ReviewService) GetAllForSelect() ([]*models.Review, error) {
	var items []*models.Review

	query := s.DB.Model(&models.Review{})

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

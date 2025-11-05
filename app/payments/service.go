package payments

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
	CreatePaymentEvent = "payments.create"
	UpdatePaymentEvent = "payments.update"
	DeletePaymentEvent = "payments.delete"
)

type PaymentService struct {
	DB      *gorm.DB
	Emitter *emitter.Emitter
	Storage *storage.ActiveStorage
	Logger  logger.Logger
}

func NewPaymentService(db *gorm.DB, emitter *emitter.Emitter, storage *storage.ActiveStorage, logger logger.Logger) *PaymentService {
	return &PaymentService{
		DB:      db,
		Logger:  logger,
		Emitter: emitter,
		Storage: storage,
	}
}

// applySorting applies sorting to the query based on the sort and order parameters
func (s *PaymentService) applySorting(query *gorm.DB, sortBy *string, sortOrder *string) {
	// Valid sortable fields for Payment
	validSortFields := map[string]string{
		"id":             "id",
		"created_at":     "created_at",
		"updated_at":     "updated_at",
		"amount":         "amount",
		"payment_method": "payment_method",
		"payment_status": "payment_status",
		"transaction_id": "transaction_id",
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

func (s *PaymentService) Create(req *models.CreatePaymentRequest) (*models.Payment, error) {
	item := &models.Payment{
		UserId:        req.UserId,
		CourseId:      req.CourseId,
		Amount:        req.Amount,
		PaymentMethod: req.PaymentMethod,
		PaymentStatus: req.PaymentStatus,
		TransactionId: req.TransactionId,
	}

	if err := s.DB.Create(item).Error; err != nil {
		s.Logger.Error("failed to create payment", logger.String("error", err.Error()))
		return nil, err
	}

	// Emit create event
	s.Emitter.Emit(CreatePaymentEvent, item)

	return s.GetById(item.Id)
}

func (s *PaymentService) Update(id uint, req *models.UpdatePaymentRequest) (*models.Payment, error) {
	item := &models.Payment{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find payment for update",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Validate request
	if err := ValidatePaymentUpdateRequest(req, id); err != nil {
		return nil, err
	}

	// Update fields directly on the model
	// For foreign key relationships
	if req.UserId != 0 {
		item.UserId = req.UserId
	}
	// For foreign key relationships
	if req.CourseId != 0 {
		item.CourseId = req.CourseId
	}
	// For non-pointer string fields
	if req.TransactionId != "" {
		item.TransactionId = req.TransactionId
	}

	if err := s.DB.Save(item).Error; err != nil {
		s.Logger.Error("failed to update payment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Handle many-to-many relationships

	result, err := s.GetById(item.Id)
	if err != nil {
		s.Logger.Error("failed to get updated payment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	// Emit update event
	s.Emitter.Emit(UpdatePaymentEvent, result)

	return result, nil
}

func (s *PaymentService) Delete(id uint) error {
	item := &models.Payment{}
	if err := s.DB.First(item, id).Error; err != nil {
		s.Logger.Error("failed to find payment for deletion",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Delete file attachments if any

	if err := s.DB.Delete(item).Error; err != nil {
		s.Logger.Error("failed to delete payment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return err
	}

	// Emit delete event
	s.Emitter.Emit(DeletePaymentEvent, item)

	return nil
}

func (s *PaymentService) GetById(id uint) (*models.Payment, error) {
	item := &models.Payment{}

	query := item.Preload(s.DB)
	if err := query.First(item, id).Error; err != nil {
		s.Logger.Error("failed to get payment",
			logger.String("error", err.Error()),
			logger.Int("id", int(id)))
		return nil, err
	}

	return item, nil
}

func (s *PaymentService) GetAll(page *int, limit *int, sortBy *string, sortOrder *string) (*types.PaginatedResponse, error) {
	var items []*models.Payment
	var total int64

	query := s.DB.Model(&models.Payment{})
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
		s.Logger.Error("failed to count payments",
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
	// query = (&models.Payment{}).Preload(query)

	// Execute query
	if err := query.Find(&items).Error; err != nil {
		s.Logger.Error("failed to get payments",
			logger.String("error", err.Error()))
		return nil, err
	}

	// Convert to response type
	responses := make([]*models.PaymentListResponse, len(items))
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
func (s *PaymentService) GetAllForSelect() ([]*models.Payment, error) {
	var items []*models.Payment

	query := s.DB.Model(&models.Payment{})

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

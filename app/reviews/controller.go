package reviews

import (
	"net/http"
	"strconv"
	"strings"

	"base/app/models"
	"base/core/router"
	"base/core/storage"
	"base/core/types"
)

type ReviewController struct {
	Service *ReviewService
	Storage *storage.ActiveStorage
}

func NewReviewController(service *ReviewService, storage *storage.ActiveStorage) *ReviewController {
	return &ReviewController{
		Service: service,
		Storage: storage,
	}
}

func (c *ReviewController) Routes(router *router.RouterGroup) {
	// Main CRUD endpoints - specific routes MUST come before parameterized routes
	router.GET("/reviews", c.List)          // Paginated list
	router.POST("/reviews", c.Create)       // Create
	router.GET("/reviews/all", c.ListAll)   // Unpaginated list - MUST be before /:id
	router.GET("/reviews/:id", c.Get)       // Get by ID - MUST be after /all
	router.PUT("/reviews/:id", c.Update)    // Update
	router.DELETE("/reviews/:id", c.Delete) // Delete

	//Upload endpoints for each file field
}

// CreateReview godoc
// @Summary Create a new Review
// @Description Create a new Review with the input payload
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param reviews body models.CreateReviewRequest true "Create Review request"
// @Success 201 {object} models.ReviewResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /reviews [post]
func (c *ReviewController) Create(ctx *router.Context) error {
	var req models.CreateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	item, err := c.Service.Create(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create item: " + err.Error()})
	}

	return ctx.JSON(http.StatusCreated, item.ToResponse())
}

// GetReview godoc
// @Summary Get a Review
// @Description Get a Review by its id
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Review id"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /reviews/{id} [get]
func (c *ReviewController) Get(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	item, err := c.Service.GetById(uint(id))
	if err != nil {
		return ctx.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Item not found"})
	}

	return ctx.JSON(http.StatusOK, item.ToResponse())
}

// ListReviews godoc
// @Summary List reviews
// @Description Get a list of reviews
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort field (id, created_at, updated_at,rating,comment,)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /reviews [get]
func (c *ReviewController) List(ctx *router.Context) error {
	var page, limit *int
	var sortBy, sortOrder *string

	// Parse page parameter
	if pageStr := ctx.Query("page"); pageStr != "" {
		if pageNum, err := strconv.Atoi(pageStr); err == nil && pageNum > 0 {
			page = &pageNum
		} else {
			return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid page number"})
		}
	}

	// Parse limit parameter
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limitNum, err := strconv.Atoi(limitStr); err == nil && limitNum > 0 {
			limit = &limitNum
		} else {
			return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid limit number"})
		}
	}

	// Parse sort parameters
	if sortStr := ctx.Query("sort"); sortStr != "" {
		sortBy = &sortStr
	}

	if orderStr := ctx.Query("order"); orderStr != "" {
		if orderStr == "asc" || orderStr == "desc" {
			sortOrder = &orderStr
		} else {
			return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid sort order. Use 'asc' or 'desc'"})
		}
	}

	paginatedResponse, err := c.Service.GetAll(page, limit, sortBy, sortOrder)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch items: " + err.Error()})
	}

	return ctx.JSON(http.StatusOK, paginatedResponse)
}

// ListAllReviews godoc
// @Summary List all reviews for select options
// @Description Get a simplified list of all reviews with id and name only (for dropdowns/select boxes)
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.ReviewSelectOption
// @Failure 500 {object} types.ErrorResponse
// @Router /reviews/all [get]
func (c *ReviewController) ListAll(ctx *router.Context) error {
	items, err := c.Service.GetAllForSelect()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch select options: " + err.Error()})
	}

	// Convert to select options
	var selectOptions []*models.ReviewSelectOption
	for _, item := range items {
		selectOptions = append(selectOptions, item.ToSelectOption())
	}

	return ctx.JSON(http.StatusOK, selectOptions)
}

// UpdateReview godoc
// @Summary Update a Review
// @Description Update a Review by its id
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Review id"
// @Param reviews body models.UpdateReviewRequest true "Update Review request"
// @Success 200 {object} models.ReviewResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /reviews/{id} [put]
func (c *ReviewController) Update(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	var req models.UpdateReviewRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	item, err := c.Service.Update(uint(id), &req)
	if err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return ctx.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Item not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to update item: " + err.Error()})
	}

	return ctx.JSON(http.StatusOK, item.ToResponse())
}

// DeleteReview godoc
// @Summary Delete a Review
// @Description Delete a Review by its id
// @Tags App/Review
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Review id"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /reviews/{id} [delete]
func (c *ReviewController) Delete(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	if err := c.Service.Delete(uint(id)); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			return ctx.JSON(http.StatusNotFound, types.ErrorResponse{Error: "Item not found"})
		}
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to delete item: " + err.Error()})
	}

	ctx.Status(http.StatusNoContent)
	return nil
}

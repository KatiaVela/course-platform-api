package course_progress_logs

import (
	"net/http"
	"strconv"
	"strings"

	"base/app/models"
	"base/core/router"
	"base/core/storage"
	"base/core/types"
)

type CourseProgressLogController struct {
	Service *CourseProgressLogService
	Storage *storage.ActiveStorage
}

func NewCourseProgressLogController(service *CourseProgressLogService, storage *storage.ActiveStorage) *CourseProgressLogController {
	return &CourseProgressLogController{
		Service: service,
		Storage: storage,
	}
}

func (c *CourseProgressLogController) Routes(router *router.RouterGroup) {
	// Main CRUD endpoints - specific routes MUST come before parameterized routes
	router.GET("/course-progress-logs", c.List)          // Paginated list
	router.POST("/course-progress-logs", c.Create)       // Create
	router.GET("/course-progress-logs/all", c.ListAll)   // Unpaginated list - MUST be before /:id
	router.GET("/course-progress-logs/:id", c.Get)       // Get by ID - MUST be after /all
	router.PUT("/course-progress-logs/:id", c.Update)    // Update
	router.DELETE("/course-progress-logs/:id", c.Delete) // Delete

	//Upload endpoints for each file field
}

// CreateCourseProgressLog godoc
// @Summary Create a new CourseProgressLog
// @Description Create a new CourseProgressLog with the input payload
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param course-progress-logs body models.CreateCourseProgressLogRequest true "Create CourseProgressLog request"
// @Success 201 {object} models.CourseProgressLogResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-progress-logs [post]
func (c *CourseProgressLogController) Create(ctx *router.Context) error {
	var req models.CreateCourseProgressLogRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	item, err := c.Service.Create(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create item: " + err.Error()})
	}

	return ctx.JSON(http.StatusCreated, item.ToResponse())
}

// GetCourseProgressLog godoc
// @Summary Get a CourseProgressLog
// @Description Get a CourseProgressLog by its id
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseProgressLog id"
// @Success 200 {object} models.CourseProgressLogResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /course-progress-logs/{id} [get]
func (c *CourseProgressLogController) Get(ctx *router.Context) error {
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

// ListCourseProgressLogs godoc
// @Summary List course-progress-logs
// @Description Get a list of course-progress-logs
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort field (id, created_at, updated_at,completed_at,)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-progress-logs [get]
func (c *CourseProgressLogController) List(ctx *router.Context) error {
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

// ListAllCourseProgressLogs godoc
// @Summary List all course-progress-logs for select options
// @Description Get a simplified list of all course-progress-logs with id and name only (for dropdowns/select boxes)
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.CourseProgressLogSelectOption
// @Failure 500 {object} types.ErrorResponse
// @Router /course-progress-logs/all [get]
func (c *CourseProgressLogController) ListAll(ctx *router.Context) error {
	items, err := c.Service.GetAllForSelect()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch select options: " + err.Error()})
	}

	// Convert to select options
	var selectOptions []*models.CourseProgressLogSelectOption
	for _, item := range items {
		selectOptions = append(selectOptions, item.ToSelectOption())
	}

	return ctx.JSON(http.StatusOK, selectOptions)
}

// UpdateCourseProgressLog godoc
// @Summary Update a CourseProgressLog
// @Description Update a CourseProgressLog by its id
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseProgressLog id"
// @Param course-progress-logs body models.UpdateCourseProgressLogRequest true "Update CourseProgressLog request"
// @Success 200 {object} models.CourseProgressLogResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-progress-logs/{id} [put]
func (c *CourseProgressLogController) Update(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	var req models.UpdateCourseProgressLogRequest
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

// DeleteCourseProgressLog godoc
// @Summary Delete a CourseProgressLog
// @Description Delete a CourseProgressLog by its id
// @Tags App/CourseProgressLog
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseProgressLog id"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-progress-logs/{id} [delete]
func (c *CourseProgressLogController) Delete(ctx *router.Context) error {
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

package course_tag_relations

import (
	"net/http"
	"strconv"
	"strings"

	"base/app/models"
	"base/core/router"
	"base/core/storage"
	"base/core/types"
)

type CourseTagRelationController struct {
	Service *CourseTagRelationService
	Storage *storage.ActiveStorage
}

func NewCourseTagRelationController(service *CourseTagRelationService, storage *storage.ActiveStorage) *CourseTagRelationController {
	return &CourseTagRelationController{
		Service: service,
		Storage: storage,
	}
}

func (c *CourseTagRelationController) Routes(router *router.RouterGroup) {
	// Main CRUD endpoints - specific routes MUST come before parameterized routes
	router.GET("/course-tag-relations", c.List)          // Paginated list
	router.POST("/course-tag-relations", c.Create)       // Create
	router.GET("/course-tag-relations/all", c.ListAll)   // Unpaginated list - MUST be before /:id
	router.GET("/course-tag-relations/:id", c.Get)       // Get by ID - MUST be after /all
	router.PUT("/course-tag-relations/:id", c.Update)    // Update
	router.DELETE("/course-tag-relations/:id", c.Delete) // Delete

	//Upload endpoints for each file field
}

// CreateCourseTagRelation godoc
// @Summary Create a new CourseTagRelation
// @Description Create a new CourseTagRelation with the input payload
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param course-tag-relations body models.CreateCourseTagRelationRequest true "Create CourseTagRelation request"
// @Success 201 {object} models.CourseTagRelationResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-tag-relations [post]
func (c *CourseTagRelationController) Create(ctx *router.Context) error {
	var req models.CreateCourseTagRelationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	item, err := c.Service.Create(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create item: " + err.Error()})
	}

	return ctx.JSON(http.StatusCreated, item.ToResponse())
}

// GetCourseTagRelation godoc
// @Summary Get a CourseTagRelation
// @Description Get a CourseTagRelation by its id
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseTagRelation id"
// @Success 200 {object} models.CourseTagRelationResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /course-tag-relations/{id} [get]
func (c *CourseTagRelationController) Get(ctx *router.Context) error {
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

// ListCourseTagRelations godoc
// @Summary List course-tag-relations
// @Description Get a list of course-tag-relations
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort field (id, created_at, updated_at,)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-tag-relations [get]
func (c *CourseTagRelationController) List(ctx *router.Context) error {
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

// ListAllCourseTagRelations godoc
// @Summary List all course-tag-relations for select options
// @Description Get a simplified list of all course-tag-relations with id and name only (for dropdowns/select boxes)
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.CourseTagRelationSelectOption
// @Failure 500 {object} types.ErrorResponse
// @Router /course-tag-relations/all [get]
func (c *CourseTagRelationController) ListAll(ctx *router.Context) error {
	items, err := c.Service.GetAllForSelect()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch select options: " + err.Error()})
	}

	// Convert to select options
	var selectOptions []*models.CourseTagRelationSelectOption
	for _, item := range items {
		selectOptions = append(selectOptions, item.ToSelectOption())
	}

	return ctx.JSON(http.StatusOK, selectOptions)
}

// UpdateCourseTagRelation godoc
// @Summary Update a CourseTagRelation
// @Description Update a CourseTagRelation by its id
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseTagRelation id"
// @Param course-tag-relations body models.UpdateCourseTagRelationRequest true "Update CourseTagRelation request"
// @Success 200 {object} models.CourseTagRelationResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-tag-relations/{id} [put]
func (c *CourseTagRelationController) Update(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	var req models.UpdateCourseTagRelationRequest
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

// DeleteCourseTagRelation godoc
// @Summary Delete a CourseTagRelation
// @Description Delete a CourseTagRelation by its id
// @Tags App/CourseTagRelation
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "CourseTagRelation id"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /course-tag-relations/{id} [delete]
func (c *CourseTagRelationController) Delete(ctx *router.Context) error {
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

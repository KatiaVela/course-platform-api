package courses

import (
	"net/http"
	"strconv"
	"strings"

	"base/app/models"
	"base/core/router"
	"base/core/storage"
	"base/core/types"
)

type CourseController struct {
	Service *CourseService
	Storage *storage.ActiveStorage
}

func NewCourseController(service *CourseService, storage *storage.ActiveStorage) *CourseController {
	return &CourseController{
		Service: service,
		Storage: storage,
	}
}

func (c *CourseController) Routes(router *router.RouterGroup) {
	// Main CRUD endpoints - specific routes MUST come before parameterized routes
	router.GET("/courses", c.List)          // Paginated list
	router.POST("/courses", c.Create)       // Create
	router.GET("/courses/all", c.ListAll)   // Unpaginated list - MUST be before /:id
	router.GET("/courses/:id", c.Get)       // Get by ID - MUST be after /all
	router.PUT("/courses/:id", c.Update)    // Update
	router.DELETE("/courses/:id", c.Delete) // Delete

	//Upload endpoints for each file field
}

// CreateCourse godoc
// @Summary Create a new Course
// @Description Create a new Course with the input payload
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param courses body models.CreateCourseRequest true "Create Course request"
// @Success 201 {object} models.CourseResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /courses [post]
func (c *CourseController) Create(ctx *router.Context) error {
	var req models.CreateCourseRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: err.Error()})
	}

	item, err := c.Service.Create(&req)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to create item: " + err.Error()})
	}

	return ctx.JSON(http.StatusCreated, item.ToResponse())
}

// GetCourse godoc
// @Summary Get a Course
// @Description Get a Course by its id
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course id"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Router /courses/{id} [get]
func (c *CourseController) Get(ctx *router.Context) error {
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

// ListCourses godoc
// @Summary List courses
// @Description Get a list of courses
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Number of items per page"
// @Param sort query string false "Sort field (id, created_at, updated_at,title,slug,description,price,level,language,thumbnail_url,status,duration,)"
// @Param order query string false "Sort order (asc, desc)"
// @Success 200 {object} types.PaginatedResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /courses [get]
func (c *CourseController) List(ctx *router.Context) error {
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

// ListAllCourses godoc
// @Summary List all courses for select options
// @Description Get a simplified list of all courses with id and name only (for dropdowns/select boxes)
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Success 200 {array} models.CourseSelectOption
// @Failure 500 {object} types.ErrorResponse
// @Router /courses/all [get]
func (c *CourseController) ListAll(ctx *router.Context) error {
	items, err := c.Service.GetAllForSelect()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, types.ErrorResponse{Error: "Failed to fetch select options: " + err.Error()})
	}

	// Convert to select options
	var selectOptions []*models.CourseSelectOption
	for _, item := range items {
		selectOptions = append(selectOptions, item.ToSelectOption())
	}

	return ctx.JSON(http.StatusOK, selectOptions)
}

// UpdateCourse godoc
// @Summary Update a Course
// @Description Update a Course by its id
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course id"
// @Param courses body models.UpdateCourseRequest true "Update Course request"
// @Success 200 {object} models.CourseResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 404 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /courses/{id} [put]
func (c *CourseController) Update(ctx *router.Context) error {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, types.ErrorResponse{Error: "Invalid id format"})
	}

	var req models.UpdateCourseRequest
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

// DeleteCourse godoc
// @Summary Delete a Course
// @Description Delete a Course by its id
// @Tags App/Course
// @Security ApiKeyAuth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course id"
// @Success 200 {object} types.SuccessResponse
// @Failure 400 {object} types.ErrorResponse
// @Failure 500 {object} types.ErrorResponse
// @Router /courses/{id} [delete]
func (c *CourseController) Delete(ctx *router.Context) error {
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

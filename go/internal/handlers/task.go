package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/org/task-api/internal/dto"
	"github.com/org/task-api/internal/exceptions"
	"github.com/org/task-api/internal/service"
)

// TaskHandler handles HTTP requests for tasks
type TaskHandler struct {
	taskService *service.TaskService
}

// NewTaskHandler creates a new TaskHandler
func NewTaskHandler(ts *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: ts}
}

// AddTask handles POST /api/task
// @Summary Create a new task
// @Description Creates a new task with the given name, description, and due date
// @Accept json
// @Produce json
// @Param task body dto.TaskRequest true "Task creation request"
// @Success 201 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/task [post]
// @Tags tasks
func (h *TaskHandler) AddTask(c *gin.Context) {
	var req dto.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Check if it's a JSON parse error (empty body or invalid JSON)
		errStr := err.Error()
		if strings.Contains(errStr, "EOF") || strings.Contains(errStr, "expecting {") {
			c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
				exceptions.NewRequestException("Request must be json!")))
			return
		}
		// For validation errors (like missing required fields), pass to service for proper validation
		// The service will return the proper error message
		resp, err := h.taskService.AddTask(req)
		if err != nil {
			h.handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
		return
	}

	resp, err := h.taskService.AddTask(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetTasks handles GET /api/task
// @Summary Get all tasks
// @Description Retrieves all tasks, optionally filtered by completion status and due date
// @Produce json
// @Param areDone query bool false "Filter by completion status"
// @Param when query string false "Filter by due date (Y-m-d)"
// @Success 200 {object} dto.TaskListResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/task [get]
// @Tags tasks
func (h *TaskHandler) GetTasks(c *gin.Context) {
	var areDone *bool
	var when *string

	if c.Query("areDone") != "" {
		b := c.Query("areDone") == "true"
		areDone = &b
	}

	if c.Query("when") != "" {
		w := c.Query("when")
		when = &w
		// Validate date format
		if _, err := time.Parse("2006-01-02", w); err != nil {
			c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
				exceptions.NewRequestException("When is not valid!")))
			return
		}
	}

	resp, err := h.taskService.GetAllTasks(areDone, when)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// GetTask handles GET /api/task/:id
// @Summary Get a specific task
// @Description Retrieves a single task by its ID
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/task/{id} [get]
// @Tags tasks
func (h *TaskHandler) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
			exceptions.NewRequestException("ID must be an integer!")))
		return
	}

	resp, err := h.taskService.GetTask(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// UpdateTask handles PUT /api/task
// @Summary Update an existing task
// @Description Updates an existing task with the provided details
// @Accept json
// @Produce json
// @Param task body dto.TaskUpdateRequest true "Task update request"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/task [put]
// @Tags tasks
func (h *TaskHandler) UpdateTask(c *gin.Context) {
	var req dto.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "EOF") || strings.Contains(errStr, "expecting {") {
			c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
				exceptions.NewRequestException("Request must be json!")))
			return
		}
		// For validation errors, pass to service for proper validation
		resp, err := h.taskService.UpdateTask(req)
		if err != nil {
			h.handleError(c, err)
			return
		}
		c.JSON(http.StatusOK, gin.H{"data": resp})
		return
	}

	resp, err := h.taskService.UpdateTask(req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// CompleteTask handles PUT /api/task/:id/complete
// @Summary Mark a task as completed
// @Description Marks a specific task as completed
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/task/{id}/complete [put]
// @Tags tasks
func (h *TaskHandler) CompleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
			exceptions.NewRequestException("ID must be an integer!")))
		return
	}

	resp, err := h.taskService.CompleteTask(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// DeleteTask handles DELETE /api/task/:id
// @Summary Delete a task
// @Description Deletes a specific task
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} dto.TaskResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /api/task/{id} [delete]
// @Tags tasks
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(
			exceptions.NewRequestException("ID must be an integer!")))
		return
	}

	resp, err := h.taskService.DeleteTask(uint(id))
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": resp})
}

// handleError handles errors and returns appropriate HTTP responses
func (h *TaskHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *exceptions.TaskNotFoundException:
		c.JSON(http.StatusNotFound, exceptions.FormatErrorResponse(e))
	case *exceptions.InvalidTaskException:
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(e))
	case *exceptions.ValidationException:
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(e))
	case *exceptions.RequestException:
		c.JSON(http.StatusBadRequest, exceptions.FormatErrorResponse(e))
	default:
		c.JSON(http.StatusInternalServerError, exceptions.FormatErrorResponse(
			exceptions.NewRequestException("Internal server error")))
	}
}
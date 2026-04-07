package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/workspace/go-api/internal/service"
)

// TaskHandler handles HTTP requests for task endpoints.
type TaskHandler struct {
	service *service.TaskService
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(svc *service.TaskService) *TaskHandler {
	return &TaskHandler{service: svc}
}

// RegisterRoutes registers all task routes on the given chi router.
// Route paths match the PHP Symfony routes exactly.
func (h *TaskHandler) RegisterRoutes(r chi.Router) {
	r.Post("/api/task", h.AddTask)
	r.Get("/api/task", h.GetTasks)
	r.Get("/api/task/{taskId}", h.GetTask)
	r.Put("/api/task", h.UpdateTask)
	r.Put("/api/task/{taskId}/complete", h.CompleteTask)
	r.Delete("/api/task/{taskId}", h.DeleteTask)
}

// AddTask handles POST /api/task
func (h *TaskHandler) AddTask(w http.ResponseWriter, r *http.Request) {
	var req CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Request must be json!")
		return
	}

	task, err := h.service.AddTask(req.Name, req.Description, req.When)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{Data: TaskToData(*task)})
}

// GetTasks handles GET /api/task
func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	var areDone *bool
	var when *time.Time

	if r.URL.Query().Has("areDone") {
		val := r.URL.Query().Get("areDone")
		b := val == "1" || val == "true"
		areDone = &b
	}

	if r.URL.Query().Has("when") {
		whenStr := r.URL.Query().Get("when")
		parsed, err := time.Parse("2006-01-02", whenStr)
		if err != nil {
			// Matches PHP: InvalidArgumentException("When is not valid!", 400)
			// PHP's ExceptionListener doesn't handle InvalidArgumentException specially,
			// so it falls through to HTTP 500.
			writeJSONError(w, http.StatusInternalServerError, "When is not valid!")
			return
		}
		when = &parsed
	}

	tasks, err := h.service.GetAllTasks(areDone, when)
	if err != nil {
		writeError(w, err)
		return
	}

	data := make([]TaskData, len(tasks))
	for i, t := range tasks {
		data[i] = TaskToData(t)
	}

	writeJSON(w, http.StatusOK, TaskListResponse{Data: data})
}

// GetTask handles GET /api/task/{taskId}
func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Task not found!")
		return
	}

	task, err := h.service.GetTask(taskID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{Data: TaskToData(*task)})
}

// UpdateTask handles PUT /api/task
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var req UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "Request must be json!")
		return
	}

	// Extract integer ID for the service call.
	// The raw req.ID (interface{}) is passed to the validator for type checking.
	var id int
	if req.ID != nil {
		if fID, ok := req.ID.(float64); ok {
			id = int(fID)
		}
	}

	task, err := h.service.UpdateTask(id, req.Name, req.Description, req.When, req.ID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{Data: TaskToData(*task)})
}

// CompleteTask handles PUT /api/task/{taskId}/complete
func (h *TaskHandler) CompleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Task not found!")
		return
	}

	task, err := h.service.CompleteTask(taskID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{Data: TaskToData(*task)})
}

// DeleteTask handles DELETE /api/task/{taskId}
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(chi.URLParam(r, "taskId"))
	if err != nil {
		writeJSONError(w, http.StatusNotFound, "Task not found!")
		return
	}

	task, err := h.service.DeleteTask(taskID)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, TaskResponse{Data: TaskToData(*task)})
}

// writeJSON writes a JSON response with the given status code.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeJSONError writes a simple error response matching PHP's ExceptionListener format.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Data: "", Error: message})
}

// writeError maps service-layer errors to HTTP responses.
// This replaces PHP's ExceptionListener functionality.
func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")

	switch e := err.(type) {
	case *service.ValidationError:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ValidationErrorResponse{
			Data:     "",
			Error:    e.Error(),
			Messages: e.Messages,
		})
	case *service.TaskNotFoundError:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Data: "", Error: e.Error()})
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Data: "", Error: err.Error()})
	}
}

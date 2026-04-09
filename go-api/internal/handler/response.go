package handler

import (
	"fmt"
	"time"

	"github.com/workspace/go-api/internal/model"
)

// TaskResponse wraps a single task for JSON output, matching PHP Fractal's Item format.
type TaskResponse struct {
	Data TaskData `json:"data"`
}

// TaskListResponse wraps a collection of tasks, matching PHP Fractal's Collection format.
type TaskListResponse struct {
	Data []TaskData `json:"data"`
}

// TaskData is the JSON-serializable representation of a task.
// Field names and structure match PHP's TaskTransformer output exactly.
type TaskData struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	When        string      `json:"when"`
	Done        bool        `json:"done"`
	CreatedAt   PHPDateTime `json:"created_at"`
	UpdatedAt   PHPDateTime `json:"updated_at"`
}

// PHPDateTime replicates PHP's \DateTime JSON serialization format:
// {"date":"2026-04-07 10:30:15.000000","timezone_type":3,"timezone":"UTC"}
type PHPDateTime struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

// NewPHPDateTime converts a Go time.Time to the PHP DateTime JSON format.
func NewPHPDateTime(t time.Time) PHPDateTime {
	return PHPDateTime{
		Date:         t.UTC().Format("2006-01-02 15:04:05.000000"),
		TimezoneType: 3,
		Timezone:     "UTC",
	}
}

// NewPHPDateTimePtr handles a nullable time pointer.
// Returns a zero-time PHPDateTime if the pointer is nil.
func NewPHPDateTimePtr(t *time.Time) PHPDateTime {
	if t == nil {
		return NewPHPDateTime(time.Time{})
	}
	return NewPHPDateTime(*t)
}

// TaskToData converts a model.Task to the JSON-ready TaskData struct.
func TaskToData(t model.Task) TaskData {
	return TaskData{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		When:        t.When.Format("2006-01-02"),
		Done:        t.Done,
		CreatedAt:   NewPHPDateTime(t.CreatedAt),
		UpdatedAt:   NewPHPDateTimePtr(t.UpdatedAt),
	}
}

// ErrorResponse is the JSON format for non-validation errors (404, 500, etc.).
// Matches PHP's ExceptionListener output: {"data":"","error":"..."}
type ErrorResponse struct {
	Data  string `json:"data"`
	Error string `json:"error"`
}

// ValidationErrorResponse is the JSON format for validation errors (400).
// Matches PHP's ExceptionListener output for ValidationException:
// {"data":"","error":"Task is not valid!","messages":["..."]}
type ValidationErrorResponse struct {
	Data     string   `json:"data"`
	Error    string   `json:"error"`
	Messages []string `json:"messages"`
}

// CreateTaskRequest represents the JSON body for POST /api/task.
type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	When        string `json:"when"`
}

// UpdateTaskRequest represents the JSON body for PUT /api/task.
// ID is interface{} because PHP accepts null, int, or string for validation.
type UpdateTaskRequest struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	When        string      `json:"when"`
}

// String returns the string representation of the UpdateTaskRequest, for debugging.
func (r UpdateTaskRequest) String() string {
	return fmt.Sprintf("UpdateTaskRequest{ID:%v, Name:%s, Description:%s, When:%s}",
		r.ID, r.Name, r.Description, r.When)
}

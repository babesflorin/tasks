package dto

import "time"

// TaskRequest represents the request body for creating a task
type TaskRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	When        string `json:"when" binding:"required"`
}

// TaskUpdateRequest represents the request body for updating a task
type TaskUpdateRequest struct {
	ID          uint   `json:"id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
	When        string `json:"when" binding:"required"`
}

// TaskResponse represents the response body for a single task
type TaskResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	When        string    `json:"when"`
	Done        bool      `json:"done"`
	CreatedAt   time.Time `json:"createdAt,omitempty"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty"`
}

// TaskListResponse represents the response body for multiple tasks
type TaskListResponse struct {
	Data []TaskResponse `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string   `json:"error,omitempty"`
	Messages []string `json:"messages,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}
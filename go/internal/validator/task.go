package validator

import (
	"errors"
	"time"

	"github.com/org/task-api/internal/dto"
)

// TaskValidator validates task requests
type TaskValidator struct{}

// NewTaskValidator creates a new TaskValidator
func NewTaskValidator() *TaskValidator {
	return &TaskValidator{}
}

// ValidationError holds validation error messages
type ValidationError struct {
	Messages []string
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return "Task is not valid!"
}

// ValidateTask validates a create task request
func (v *TaskValidator) ValidateTask(req dto.TaskRequest) *ValidationError {
	var messages []string

	if req.Name == "" {
		messages = append(messages, "Task name is not valid!")
	}

	if req.Description == "" {
		messages = append(messages, "Task description is not valid!")
	}

	if req.When == "" {
		messages = append(messages, "Task must have a date!")
	} else if _, err := time.Parse("2006-01-02", req.When); err != nil {
		messages = append(messages, "`when` is not a valid date!")
	}
	// Note: Future date validation is handled in service layer

	if len(messages) > 0 {
		return &ValidationError{Messages: messages}
	}

	return nil
}

// ValidateTaskUpdate validates an update task request
func (v *TaskValidator) ValidateTaskUpdate(req dto.TaskUpdateRequest) *ValidationError {
	var messages []string

	if req.ID == 0 {
		messages = append(messages, "We need an id to know which entity to update!")
	}

	if req.Name == "" {
		messages = append(messages, "Task name is not valid!")
	}

	if req.Description == "" {
		messages = append(messages, "Task description is not valid!")
	}

	if req.When == "" {
		messages = append(messages, "Task must have a date!")
	} else if _, err := time.Parse("2006-01-02", req.When); err != nil {
		messages = append(messages, "`when` is not a valid date!")
	}
	// Note: Future date validation is handled in service layer

	if len(messages) > 0 {
		return &ValidationError{Messages: messages}
	}

	return nil
}

// IsValidDate checks if a date string is valid Y-m-d format
func IsValidDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// ValidateDate validates that date is valid and not in the past
func ValidateDate(dateStr string) error {
	if dateStr == "" {
		return errors.New("Task must have a date!")
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return errors.New("`when` is not a valid date!")
	}

	today := time.Now().Truncate(24 * time.Hour)
	if t.Before(today) {
		return errors.New("You can't do a task in a the past. Or you can?")
	}

	return nil
}
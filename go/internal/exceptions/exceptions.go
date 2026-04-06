package exceptions

// ValidationException represents a general validation error
type ValidationException struct {
	Message string
}

func (e *ValidationException) Error() string {
	return e.Message
}

// NewValidationException creates a new ValidationException
func NewValidationException(message string) *ValidationException {
	return &ValidationException{Message: message}
}

// InvalidTaskException represents an invalid task error
type InvalidTaskException struct {
	Messages []string
}

func (e *InvalidTaskException) Error() string {
	return "Task is not valid!"
}

// NewInvalidTaskException creates a new InvalidTaskException
func NewInvalidTaskException(messages []string) *InvalidTaskException {
	return &InvalidTaskException{Messages: messages}
}

// TaskNotFoundException represents a task not found error
type TaskNotFoundException struct {
	Message string
}

func (e *TaskNotFoundException) Error() string {
	if e.Message == "" {
		return "Task not found!"
	}
	return e.Message
}

// NewTaskNotFoundException creates a new TaskNotFoundException
func NewTaskNotFoundException(message string) *TaskNotFoundException {
	return &TaskNotFoundException{Message: message}
}

// CouldNotDeleteException represents a delete error
type CouldNotDeleteException struct {
	Message string
}

func (e *CouldNotDeleteException) Error() string {
	if e.Message == "" {
		return "Could not delete the task!"
	}
	return e.Message
}

// NewCouldNotDeleteException creates a new CouldNotDeleteException
func NewCouldNotDeleteException(message string) *CouldNotDeleteException {
	return &CouldNotDeleteException{Message: message}
}

// RequestException represents an invalid request error
type RequestException struct {
	Message string
}

func (e *RequestException) Error() string {
	return e.Message
}

// NewRequestException creates a new RequestException
func NewRequestException(message string) *RequestException {
	return &RequestException{Message: message}
}

// FormatErrorResponse formats an error response
func FormatErrorResponse(err error) map[string]interface{} {
	resp := map[string]interface{}{
		"error": err.Error(),
		"data":  nil,
	}

	switch v := err.(type) {
	case *InvalidTaskException:
		resp["messages"] = v.Messages
	case *ValidationException:
		resp["messages"] = []string{v.Message}
	}

	return resp
}
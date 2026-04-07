package service

import (
	"fmt"
	"time"

	"github.com/workspace/go-api/internal/model"
	"github.com/workspace/go-api/internal/repository"
)

// TaskNotFoundError is returned when a task cannot be found by ID.
// Error message matches PHP's TaskNotFoundException exactly.
type TaskNotFoundError struct{}

func (e *TaskNotFoundError) Error() string {
	return "Task not found!"
}

// CouldNotDeleteError is returned when a task deletion fails.
// Error message wraps the underlying error, matching PHP's CouldNotDeleteException.
type CouldNotDeleteError struct {
	Msg string
}

func (e *CouldNotDeleteError) Error() string {
	return e.Msg
}

// TaskService contains the business logic for task operations.
type TaskService struct {
	repo *repository.TaskRepository
}

// NewTaskService creates a new TaskService.
func NewTaskService(repo *repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// AddTask validates and creates a new task.
func (s *TaskService) AddTask(name, description, when string) (*model.Task, error) {
	if err := ValidateTask(name, description, when, false, nil); err != nil {
		return nil, err
	}

	parsedWhen, _ := time.Parse("2006-01-02", when)
	task := &model.Task{
		Name:        name,
		Description: description,
		When:        parsedWhen,
		Done:        false,
	}

	return s.repo.SaveTask(task)
}

// GetAllTasks retrieves all tasks, optionally filtered by completion status and due date.
func (s *TaskService) GetAllTasks(areDone *bool, when *time.Time) ([]model.Task, error) {
	return s.repo.GetTasks(areDone, when)
}

// GetTask retrieves a single task by ID.
func (s *TaskService) GetTask(taskID int) (*model.Task, error) {
	task, err := s.repo.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, &TaskNotFoundError{}
	}
	return task, nil
}

// CompleteTask marks a task as done.
func (s *TaskService) CompleteTask(taskID int) (*model.Task, error) {
	task, err := s.repo.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, &TaskNotFoundError{}
	}

	task.Done = true
	return s.repo.SaveTask(task)
}

// UpdateTask validates and updates an existing task.
func (s *TaskService) UpdateTask(id int, name, description, when string, rawID interface{}) (*model.Task, error) {
	if err := ValidateTask(name, description, when, true, rawID); err != nil {
		return nil, err
	}

	task, err := s.repo.FindTaskByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, &TaskNotFoundError{}
	}

	parsedWhen, _ := time.Parse("2006-01-02", when)
	task.Name = name
	task.Description = description
	task.When = parsedWhen

	return s.repo.SaveTask(task)
}

// DeleteTask removes a task by ID, returning the task data before deletion.
func (s *TaskService) DeleteTask(taskID int) (*model.Task, error) {
	task, err := s.repo.FindTaskByID(taskID)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, &TaskNotFoundError{}
	}

	// Copy task data before deletion (matches PHP behavior)
	taskCopy := *task

	if err := s.repo.DeleteTask(task); err != nil {
		return nil, &CouldNotDeleteError{Msg: fmt.Sprintf("Could not delete the task: %s", err.Error())}
	}

	return &taskCopy, nil
}

package service

import (
	"time"

	"github.com/org/task-api/internal/dto"
	"github.com/org/task-api/internal/exceptions"
	"github.com/org/task-api/internal/models"
	"github.com/org/task-api/internal/repository"
	"github.com/org/task-api/internal/validator"
)

// TaskService handles task business logic
type TaskService struct {
	repo      *repository.TaskRepository
	validator *validator.TaskValidator
}

// NewTaskService creates a new TaskService
func NewTaskService(repo *repository.TaskRepository, v *validator.TaskValidator) *TaskService {
	return &TaskService{
		repo:      repo,
		validator: v,
	}
}

// AddTask creates a new task
func (s *TaskService) AddTask(req dto.TaskRequest) (*dto.TaskResponse, error) {
	// Validate request
	if err := s.validator.ValidateTask(req); err != nil {
		return nil, exceptions.NewInvalidTaskException(err.Messages)
	}

	// Parse date
	t, err := time.Parse("2006-01-02", req.When)
	if err != nil {
		return nil, exceptions.NewInvalidTaskException([]string{"`when` is not a valid date!"})
	}

	// Create task
	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
		When:        t,
		Done:        false,
	}

	// Save to database
	savedTask, err := s.repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	return s.toResponse(savedTask), nil
}

// GetAllTasks retrieves all tasks with optional filters
func (s *TaskService) GetAllTasks(areDone *bool, when *string) (*dto.TaskListResponse, error) {
	tasks, err := s.repo.GetTasks(areDone, when)
	if err != nil {
		return nil, err
	}

	response := &dto.TaskListResponse{
		Data: make([]dto.TaskResponse, 0, len(tasks)),
	}

	for _, task := range tasks {
		response.Data = append(response.Data, *s.toResponse(&task))
	}

	return response, nil
}

// GetTask retrieves a single task by ID
func (s *TaskService) GetTask(id uint) (*dto.TaskResponse, error) {
	task, err := s.repo.FindTaskById(id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, exceptions.NewTaskNotFoundException("Task not found!")
	}

	return s.toResponse(task), nil
}

// UpdateTask updates an existing task
func (s *TaskService) UpdateTask(req dto.TaskUpdateRequest) (*dto.TaskResponse, error) {
	// Validate request
	if err := s.validator.ValidateTaskUpdate(req); err != nil {
		return nil, exceptions.NewInvalidTaskException(err.Messages)
	}

	// Find existing task
	task, err := s.repo.FindTaskById(req.ID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, exceptions.NewTaskNotFoundException("Task not found!")
	}

	// Parse date
	t, err := time.Parse("2006-01-02", req.When)
	if err != nil {
		return nil, exceptions.NewInvalidTaskException([]string{"`when` is not a valid date!"})
	}

	// Update task
	task.Name = req.Name
	task.Description = req.Description
	task.When = t

	// Save to database
	savedTask, err := s.repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	return s.toResponse(savedTask), nil
}

// CompleteTask marks a task as completed
func (s *TaskService) CompleteTask(id uint) (*dto.TaskResponse, error) {
	task, err := s.repo.FindTaskById(id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, exceptions.NewTaskNotFoundException("Task not found!")
	}

	// Mark as complete
	task.Done = true

	// Save to database
	savedTask, err := s.repo.SaveTask(task)
	if err != nil {
		return nil, err
	}

	return s.toResponse(savedTask), nil
}

// DeleteTask deletes a task
func (s *TaskService) DeleteTask(id uint) (*dto.TaskResponse, error) {
	task, err := s.repo.FindTaskById(id)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, exceptions.NewTaskNotFoundException("Task not found!")
	}

	// Get response before deleting
	response := s.toResponse(task)

	// Delete from database
	err = s.repo.DeleteTask(task)
	if err != nil {
		return nil, exceptions.NewCouldNotDeleteException(err.Error())
	}

	return response, nil
}

// toResponse converts a Task model to TaskResponse DTO
func (s *TaskService) toResponse(task *models.Task) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:          task.ID,
		Name:        task.Name,
		Description: task.Description,
		When:        task.When.Format("2006-01-02"),
		Done:        task.Done,
		CreatedAt:   task.CreatedAt,
		UpdatedAt:   task.UpdatedAt,
	}
}
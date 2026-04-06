package repository

import (
	"time"

	"github.com/org/task-api/internal/models"
	"gorm.io/gorm"
)

// TaskRepository handles database operations for tasks
type TaskRepository struct {
	db *gorm.DB
}

// NewTaskRepository creates a new TaskRepository
func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// SaveTask saves a task to the database
func (r *TaskRepository) SaveTask(task *models.Task) (*models.Task, error) {
	task.UpdatedAt = time.Now()
	if task.ID == 0 {
		task.CreatedAt = time.Now()
	}

	err := r.db.Save(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}

// GetTasks retrieves all tasks with optional filters
func (r *TaskRepository) GetTasks(areDone *bool, when *string) ([]models.Task, error) {
	var tasks []models.Task
	query := r.db.Model(&models.Task{})

	if areDone != nil {
		query = query.Where("done = ?", *areDone)
	}

	if when != nil && *when != "" {
		t, err := time.Parse("2006-01-02", *when)
		if err == nil {
			query = query.Where("`when` = ?", t)
		}
	}

	err := query.Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

// FindTaskById finds a task by ID
func (r *TaskRepository) FindTaskById(id uint) (*models.Task, error) {
	var task models.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

// DeleteTask deletes a task
func (r *TaskRepository) DeleteTask(task *models.Task) error {
	return r.db.Delete(task).Error
}

// DB returns the database connection
func (r *TaskRepository) DB() *gorm.DB {
	return r.db
}
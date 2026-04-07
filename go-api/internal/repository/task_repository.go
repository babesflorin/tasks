package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/workspace/go-api/internal/model"
)

// TaskRepository handles all database operations for the task entity.
type TaskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository creates a new TaskRepository.
func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// SaveTask inserts a new task or updates an existing one.
// If task.ID == 0, it performs an INSERT; otherwise an UPDATE.
func (r *TaskRepository) SaveTask(t *model.Task) (*model.Task, error) {
	now := time.Now().UTC()
	t.UpdatedAt = &now

	if t.ID == 0 {
		// INSERT
		t.CreatedAt = now
		result, err := r.db.Exec(
			"INSERT INTO task (name, description, `when`, done, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)",
			t.Name, t.Description, t.When, t.Done, t.CreatedAt, t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("inserting task: %w", err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return nil, fmt.Errorf("getting last insert id: %w", err)
		}
		t.ID = int(id)
	} else {
		// UPDATE
		_, err := r.db.Exec(
			"UPDATE task SET name=?, description=?, `when`=?, done=?, updated_at=? WHERE id=?",
			t.Name, t.Description, t.When, t.Done, t.UpdatedAt, t.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("updating task %d: %w", t.ID, err)
		}
	}

	return t, nil
}

// GetTasks retrieves tasks with optional filters for completion status and due date.
func (r *TaskRepository) GetTasks(areDone *bool, when *time.Time) ([]model.Task, error) {
	query := "SELECT id, name, description, `when`, done, created_at, updated_at FROM task"
	args := []interface{}{}
	conditions := []string{}

	if areDone != nil {
		conditions = append(conditions, "done = ?")
		args = append(args, *areDone)
	}
	if when != nil {
		conditions = append(conditions, "DATE(`when`) = ?")
		args = append(args, when.Format("2006-01-02"))
	}

	if len(conditions) > 0 {
		query += " WHERE "
		for i, c := range conditions {
			if i > 0 {
				query += " AND "
			}
			query += c
		}
	}

	var tasks []model.Task
	if err := r.db.Select(&tasks, query, args...); err != nil {
		return nil, fmt.Errorf("querying tasks: %w", err)
	}

	if tasks == nil {
		tasks = []model.Task{}
	}

	return tasks, nil
}

// FindTaskByID retrieves a single task by its ID.
// Returns nil, nil if the task is not found.
func (r *TaskRepository) FindTaskByID(taskID int) (*model.Task, error) {
	var task model.Task
	err := r.db.Get(&task, "SELECT id, name, description, `when`, done, created_at, updated_at FROM task WHERE id = ?", taskID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("finding task %d: %w", taskID, err)
	}
	return &task, nil
}

// DeleteTask removes a task from the database.
func (r *TaskRepository) DeleteTask(t *model.Task) error {
	_, err := r.db.Exec("DELETE FROM task WHERE id = ?", t.ID)
	if err != nil {
		return fmt.Errorf("deleting task %d: %w", t.ID, err)
	}
	return nil
}

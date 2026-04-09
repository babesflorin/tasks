package model

import "time"

// Task represents the domain model for a task, mapped to the `task` MySQL table.
type Task struct {
	ID          int        `db:"id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	When        time.Time  `db:"when"`
	Done        bool       `db:"done"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

package models

import (
	"time"
)

// Task represents the task entity in the database
type Task struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description string    `gorm:"type:text;not null" json:"description"`
	When        time.Time `gorm:"type:date;not null" json:"when"`
	Done        bool      `gorm:"default:false" json:"done"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName specifies the table name for Task model
func (Task) TableName() string {
	return "task"
}
package models

import "time"

type Subtask struct {
	ID          string    `db:"id"`
	TaskID      string    `db:"task_id"`
	Title       string    `db:"title"`
	IsCompleted bool      `db:"is_completed"`
	CreatedAt   time.Time `db:"created_at"`
}

package models

import "time"

type Board struct {
	ID          string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Color       string    `db:"color"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type Column struct {
	ID        string    `db:"id"`
	BoardID   string    `db:"board_id"`
	Name      string    `db:"name"`
	Position  float64   `db:"position"`
	Color     string    `db:"color"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Checklist struct {
	ID       string          `db:"id"`
	TaskID   string          `db:"task_id"`
	Title    string          `db:"title"`
	Position float64         `db:"position"`
	Items    []ChecklistItem `db:"-"`
}

type ChecklistItem struct {
	ID          string  `db:"id"`
	ChecklistID string  `db:"checklist_id"`
	Title       string  `db:"title"`
	Completed   bool    `db:"completed"`
	Position    float64 `db:"position"`
}

type Label struct {
	ID    string `db:"id"`
	Name  string `db:"name"`
	Color string `db:"color"`
}

type Attachment struct {
	ID        string    `db:"id"`
	TaskID    string    `db:"task_id"`
	Filename  string    `db:"filename"`
	URL       string    `db:"url"`
	Size      int64     `db:"size"`
	CreatedAt time.Time `db:"created_at"`
}

type Comment struct {
	ID        string    `db:"id"`
	TaskID    string    `db:"task_id"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type ActivityLog struct {
	ID        string    `db:"id"`
	TaskID    *string   `db:"task_id"`
	Type      string    `db:"type"`
	OldValue  *string   `db:"old_value"`
	NewValue  *string   `db:"new_value"`
	CreatedAt time.Time `db:"created_at"`
}

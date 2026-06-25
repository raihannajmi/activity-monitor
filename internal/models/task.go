package models

import "time"

type Priority string
type Status string

const (
	PriorityLow    Priority = "low"
	PriorityMedium Priority = "medium"
	PriorityHigh   Priority = "high"

	StatusTodo       Status = "todo"
	StatusInProgress Status = "in_progress"
	StatusDone       Status = "done"
)

type Task struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Description string     `db:"description"`
	Priority    Priority   `db:"priority"`
	Status      Status     `db:"status"`
	Deadline    *time.Time `db:"deadline"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`

	// Populated on demand
	Subtasks  []Subtask  `db:"-"`
	Reminders []Reminder `db:"-"`
}

func (t *Task) SubtaskProgress() (int, int) {
	total := len(t.Subtasks)
	done := 0
	for _, s := range t.Subtasks {
		if s.IsCompleted {
			done++
		}
	}
	return done, total
}

func (t *Task) SubtaskPercent() int {
	done, total := t.SubtaskProgress()
	if total == 0 {
		return 0
	}
	return (done * 100) / total
}

func (t *Task) PriorityLabel() string {
	switch t.Priority {
	case PriorityHigh:
		return "High"
	case PriorityMedium:
		return "Medium"
	case PriorityLow:
		return "Low"
	default:
		return "Medium"
	}
}

func (t *Task) PriorityColor() string {
	switch t.Priority {
	case PriorityHigh:
		return "priority-high"
	case PriorityMedium:
		return "priority-medium"
	case PriorityLow:
		return "priority-low"
	default:
		return "priority-medium"
	}
}

func (t *Task) IsOverdue() bool {
	if t.Deadline == nil || t.Status == StatusDone {
		return false
	}
	return time.Now().After(*t.Deadline)
}

func (t *Task) StatusLabel() string {
	switch t.Status {
	case StatusTodo:
		return "Todo"
	case StatusInProgress:
		return "In Progress"
	case StatusDone:
		return "Done"
	default:
		return "Todo"
	}
}

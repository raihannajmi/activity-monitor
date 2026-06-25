package models

import "time"

type ActivityType string

const (
	ActivityTaskCreated    ActivityType = "task_created"
	ActivityTaskUpdated    ActivityType = "task_updated"
	ActivityTaskDone       ActivityType = "task_done"
	ActivitySubtaskDone    ActivityType = "subtask_done"
	ActivityReminderCreated ActivityType = "reminder_created"
	ActivityNoteCreated    ActivityType = "note_created"
	ActivityNoteUpdated    ActivityType = "note_updated"
)

type Activity struct {
	ID          string       `db:"id"`
	Type        ActivityType `db:"type"`
	ReferenceID string       `db:"reference_id"`
	Description string       `db:"description"`
	CreatedAt   time.Time    `db:"created_at"`
}

func (a *Activity) Icon() string {
	switch a.Type {
	case ActivityTaskCreated:
		return "✦"
	case ActivityTaskUpdated:
		return "✎"
	case ActivityTaskDone:
		return "✓"
	case ActivitySubtaskDone:
		return "◎"
	case ActivityReminderCreated:
		return "◷"
	case ActivityNoteCreated:
		return "✐"
	case ActivityNoteUpdated:
		return "✏"
	default:
		return "·"
	}
}

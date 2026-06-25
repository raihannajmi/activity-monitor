package models

import "time"

type Reminder struct {
	ID        string    `db:"id"`
	TaskID    string    `db:"task_id"`
	RemindAt  time.Time `db:"remind_at"`
	Note      string    `db:"note"`
	IsDone    bool      `db:"is_done"`
	CreatedAt time.Time `db:"created_at"`

	// Populated on demand
	TaskTitle string `db:"-"`
}

func (r *Reminder) IsToday() bool {
	now := time.Now()
	return r.RemindAt.Year() == now.Year() &&
		r.RemindAt.Month() == now.Month() &&
		r.RemindAt.Day() == now.Day()
}

func (r *Reminder) IsOverdue() bool {
	return !r.IsDone && time.Now().After(r.RemindAt)
}

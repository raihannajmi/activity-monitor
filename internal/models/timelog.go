package models

import "time"

type TimeLog struct {
	ID              string     `json:"id" db:"id"`
	TaskID          *string    `json:"task_id" db:"task_id"`
	SubtaskID       *string    `json:"subtask_id" db:"subtask_id"`
	StartTime       time.Time  `json:"start_time" db:"start_time"`
	EndTime         *time.Time `json:"end_time" db:"end_time"`
	DurationSeconds int        `json:"duration_seconds" db:"duration_seconds"`
	SessionType     string     `json:"session_type" db:"session_type"` // 'pomodoro', 'short_break', 'long_break', 'stopwatch'
}

type DailyRecap struct {
	PomodorosCompleted int
	TotalFocusSeconds  int
}

type TimeLogWithDetails struct {
	TimeLog
	TaskTitle    string `db:"task_title"`
	SubtaskTitle string `db:"subtask_title"`
}

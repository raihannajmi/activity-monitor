package models

import "time"

type TimeLog struct {
	ID              string     `json:"id"`
	TaskID          *string    `json:"task_id"`
	SubtaskID       *string    `json:"subtask_id"`
	StartTime       time.Time  `json:"start_time"`
	EndTime         *time.Time `json:"end_time"`
	DurationSeconds int        `json:"duration_seconds"`
	SessionType     string     `json:"session_type"` // 'pomodoro', 'short_break', 'long_break', 'stopwatch'
}

type DailyRecap struct {
	PomodorosCompleted int
	TotalFocusSeconds  int
}

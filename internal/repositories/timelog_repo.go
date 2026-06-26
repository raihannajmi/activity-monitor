package repositories

import (
	"database/sql"
	"time"

	"activity-monitor/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type TimeLogRepository struct {
	db *sqlx.DB
}

func NewTimeLogRepository(db *sqlx.DB) *TimeLogRepository {
	return &TimeLogRepository{db: db}
}

func (r *TimeLogRepository) StartSession(taskID, subtaskID *string, sessionType string) (*models.TimeLog, error) {
	log := &models.TimeLog{
		ID:          uuid.New().String(),
		TaskID:      taskID,
		SubtaskID:   subtaskID,
		StartTime:   time.Now(),
		SessionType: sessionType,
	}

	query := `
		INSERT INTO time_logs (id, task_id, subtask_id, start_time, session_type)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, log.ID, log.TaskID, log.SubtaskID, log.StartTime, log.SessionType)
	return log, err
}

func (r *TimeLogRepository) StopSession(id string, durationSeconds int) error {
	query := `
		UPDATE time_logs 
		SET end_time = ?, duration_seconds = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), durationSeconds, id)
	return err
}

func (r *TimeLogRepository) GetActiveSession() (*models.TimeLog, error) {
	query := `
		SELECT id, task_id, subtask_id, start_time, session_type
		FROM time_logs
		WHERE end_time IS NULL
		ORDER BY start_time DESC LIMIT 1
	`
	var log models.TimeLog
	err := r.db.QueryRow(query).Scan(&log.ID, &log.TaskID, &log.SubtaskID, &log.StartTime, &log.SessionType)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &log, err
}

func (r *TimeLogRepository) GetTotalTimeForTask(taskID string) (int, error) {
	query := `SELECT COALESCE(SUM(duration_seconds), 0) FROM time_logs WHERE task_id = ?`
	var total int
	err := r.db.QueryRow(query, taskID).Scan(&total)
	return total, err
}

func (r *TimeLogRepository) GetTotalTimeForSubtask(subtaskID string) (int, error) {
	query := `SELECT COALESCE(SUM(duration_seconds), 0) FROM time_logs WHERE subtask_id = ?`
	var total int
	err := r.db.QueryRow(query, subtaskID).Scan(&total)
	return total, err
}

func (r *TimeLogRepository) GetDailyRecap(date time.Time) (*models.DailyRecap, error) {
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
		SELECT 
			COALESCE(SUM(CASE WHEN session_type = 'pomodoro' AND duration_seconds >= 1500 THEN 1 ELSE 0 END), 0) as pomodoros,
			COALESCE(SUM(duration_seconds), 0) as total_seconds
		FROM time_logs
		WHERE start_time >= ? AND start_time < ?
	`
	var recap models.DailyRecap
	err := r.db.QueryRow(query, startOfDay, endOfDay).Scan(&recap.PomodorosCompleted, &recap.TotalFocusSeconds)
	return &recap, err
}

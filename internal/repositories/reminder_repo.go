package repositories

import (
	"time"

	"activity-monitor/internal/models"

	"github.com/jmoiron/sqlx"
)

type ReminderRepository struct {
	db *sqlx.DB
}

func NewReminderRepository(db *sqlx.DB) *ReminderRepository {
	return &ReminderRepository{db: db}
}

func (r *ReminderRepository) ListToday() ([]models.Reminder, error) {
	var reminders []models.Reminder
	today := time.Now().Format("2006-01-02")
	err := r.db.Select(&reminders, `
		SELECT r.id, r.task_id, r.remind_at, r.note, r.is_done, r.created_at
		FROM reminders r
		WHERE DATE(r.remind_at) = ? AND r.is_done = 0
		ORDER BY r.remind_at ASC
	`, today)
	return reminders, err
}

func (r *ReminderRepository) ListUpcoming() ([]models.Reminder, error) {
	var reminders []models.Reminder
	now := time.Now()
	err := r.db.Select(&reminders, `
		SELECT id, task_id, remind_at, note, is_done, created_at
		FROM reminders
		WHERE remind_at > ? AND is_done = 0
		ORDER BY remind_at ASC
		LIMIT 10
	`, now)
	return reminders, err
}

func (r *ReminderRepository) ListByTaskID(taskID string) ([]models.Reminder, error) {
	var reminders []models.Reminder
	err := r.db.Select(&reminders, `
		SELECT id, task_id, remind_at, note, is_done, created_at
		FROM reminders WHERE task_id = ? ORDER BY remind_at ASC
	`, taskID)
	return reminders, err
}

func (r *ReminderRepository) Create(rem *models.Reminder) error {
	_, err := r.db.Exec(`
		INSERT INTO reminders (id, task_id, remind_at, note, is_done, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`, rem.ID, rem.TaskID, rem.RemindAt, rem.Note, rem.IsDone, rem.CreatedAt)
	return err
}

func (r *ReminderRepository) MarkDone(id string) error {
	_, err := r.db.Exec(`UPDATE reminders SET is_done=1 WHERE id=?`, id)
	return err
}

func (r *ReminderRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM reminders WHERE id=?`, id)
	return err
}

func (r *ReminderRepository) CountToday() (int, error) {
	var count int
	today := time.Now().Format("2006-01-02")
	err := r.db.Get(&count, `
		SELECT COUNT(*) FROM reminders WHERE DATE(remind_at)=? AND is_done=0
	`, today)
	return count, err
}

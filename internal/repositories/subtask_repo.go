package repositories

import (
	"activity-monitor/internal/models"

	"github.com/jmoiron/sqlx"
)

type SubtaskRepository struct {
	db *sqlx.DB
}

func NewSubtaskRepository(db *sqlx.DB) *SubtaskRepository {
	return &SubtaskRepository{db: db}
}

func (r *SubtaskRepository) ListByTaskID(taskID string) ([]models.Subtask, error) {
	var subs []models.Subtask
	query := `SELECT id, task_id, title, is_completed, deadline, created_at FROM subtasks WHERE task_id = ? ORDER BY created_at ASC`
	err := r.db.Select(&subs, query, taskID)
	return subs, err
}

func (r *SubtaskRepository) Create(sub *models.Subtask) error {
	query := `
		INSERT INTO subtasks (id, task_id, title, is_completed, deadline, created_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, sub.ID, sub.TaskID, sub.Title, sub.IsCompleted, sub.Deadline, sub.CreatedAt)
	return err
}

func (r *SubtaskRepository) ToggleComplete(id string) (bool, error) {
	var current bool
	if err := r.db.Get(&current, `SELECT is_completed FROM subtasks WHERE id=?`, id); err != nil {
		return false, err
	}
	newVal := !current
	_, err := r.db.Exec(`UPDATE subtasks SET is_completed=? WHERE id=?`, newVal, id)
	return newVal, err
}

func (r *SubtaskRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM subtasks WHERE id=?`, id)
	return err
}

func (r *SubtaskRepository) AllCompleted(taskID string) (bool, error) {
	var total, done int
	if err := r.db.Get(&total, `SELECT COUNT(*) FROM subtasks WHERE task_id=?`, taskID); err != nil {
		return false, err
	}
	if total == 0 {
		return false, nil
	}
	if err := r.db.Get(&done, `SELECT COUNT(*) FROM subtasks WHERE task_id=? AND is_completed=1`, taskID); err != nil {
		return false, err
	}
	return total == done, nil
}

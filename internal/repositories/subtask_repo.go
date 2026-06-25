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
	var subtasks []models.Subtask
	err := r.db.Select(&subtasks, `
		SELECT id, task_id, title, is_completed, created_at
		FROM subtasks WHERE task_id = ? ORDER BY created_at ASC
	`, taskID)
	return subtasks, err
}

func (r *SubtaskRepository) Create(s *models.Subtask) error {
	_, err := r.db.Exec(`
		INSERT INTO subtasks (id, task_id, title, is_completed, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, s.ID, s.TaskID, s.Title, s.IsCompleted, s.CreatedAt)
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

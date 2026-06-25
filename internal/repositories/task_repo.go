package repositories

import (
	"fmt"
	"time"

	"activity-monitor/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) ListAll() ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Select(&tasks, `
		SELECT id, title, description, priority, status, deadline, created_at, updated_at
		FROM tasks
		ORDER BY
			CASE status WHEN 'done' THEN 1 ELSE 0 END,
			CASE priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END,
			created_at DESC
	`)
	return tasks, err
}

func (r *TaskRepository) ListByStatus(status string) ([]models.Task, error) {
	var tasks []models.Task
	err := r.db.Select(&tasks, `
		SELECT id, title, description, priority, status, deadline, created_at, updated_at
		FROM tasks WHERE status = ?
		ORDER BY CASE priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END, created_at DESC
	`, status)
	return tasks, err
}

func (r *TaskRepository) ListDueToday() ([]models.Task, error) {
	var tasks []models.Task
	today := time.Now().Format("2006-01-02")
	err := r.db.Select(&tasks, `
		SELECT id, title, description, priority, status, deadline, created_at, updated_at
		FROM tasks
		WHERE DATE(deadline) = ? AND status != 'done'
		ORDER BY CASE priority WHEN 'high' THEN 0 WHEN 'medium' THEN 1 ELSE 2 END
	`, today)
	return tasks, err
}

func (r *TaskRepository) CountByStatus(status string) (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM tasks WHERE status = ?`, status)
	return count, err
}

func (r *TaskRepository) CountActive() (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM tasks WHERE status != 'done'`)
	return count, err
}

func (r *TaskRepository) GetByID(id string) (*models.Task, error) {
	var task models.Task
	err := r.db.Get(&task, `
		SELECT id, title, description, priority, status, deadline, created_at, updated_at
		FROM tasks WHERE id = ?
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get task %s: %w", id, err)
	}
	return &task, nil
}

func (r *TaskRepository) Create(task *models.Task) error {
	_, err := r.db.Exec(`
		INSERT INTO tasks (id, title, description, priority, status, deadline, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`, task.ID, task.Title, task.Description, task.Priority, task.Status,
		task.Deadline, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *TaskRepository) Update(task *models.Task) error {
	task.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE tasks SET title=?, description=?, priority=?, status=?, deadline=?, updated_at=?
		WHERE id=?
	`, task.Title, task.Description, task.Priority, task.Status, task.Deadline, task.UpdatedAt, task.ID)
	return err
}

func (r *TaskRepository) UpdateStatus(id string, status models.Status) error {
	_, err := r.db.Exec(`
		UPDATE tasks SET status=?, updated_at=? WHERE id=?
	`, status, time.Now(), id)
	return err
}

func (r *TaskRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM tasks WHERE id=?`, id)
	return err
}

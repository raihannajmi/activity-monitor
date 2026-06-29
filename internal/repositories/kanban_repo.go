package repositories

import (
	"activity-monitor/internal/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type KanbanRepository struct {
	db *sqlx.DB
}

func NewKanbanRepository(db *sqlx.DB) *KanbanRepository {
	return &KanbanRepository{db: db}
}

// Boards
func (r *KanbanRepository) ListBoards() ([]models.Board, error) {
	var list []models.Board
	err := r.db.Select(&list, "SELECT id, name, description, color, created_at, updated_at FROM boards ORDER BY name")
	return list, err
}

func (r *KanbanRepository) GetBoard(id string) (*models.Board, error) {
	var b models.Board
	err := r.db.Get(&b, "SELECT id, name, description, color, created_at, updated_at FROM boards WHERE id = ?", id)
	return &b, err
}

func (r *KanbanRepository) CreateBoard(b *models.Board) error {
	_, err := r.db.Exec("INSERT INTO boards (id, name, description, color) VALUES (?, ?, ?, ?)", b.ID, b.Name, b.Description, b.Color)
	return err
}

// Columns
func (r *KanbanRepository) ListColumns(boardID string) ([]models.Column, error) {
	var list []models.Column
	err := r.db.Select(&list, "SELECT id, board_id, name, position, color, created_at, updated_at FROM columns WHERE board_id = ? ORDER BY position ASC", boardID)
	return list, err
}

func (r *KanbanRepository) CreateColumn(c *models.Column) error {
	_, err := r.db.Exec("INSERT INTO columns (id, board_id, name, position, color) VALUES (?, ?, ?, ?, ?)", c.ID, c.BoardID, c.Name, c.Position, c.Color)
	return err
}

func (r *KanbanRepository) UpdateColumn(c *models.Column) error {
	_, err := r.db.Exec("UPDATE columns SET name = ?, position = ?, color = ?, updated_at = ? WHERE id = ?", c.Name, c.Position, c.Color, time.Now(), c.ID)
	return err
}

func (r *KanbanRepository) DeleteColumn(id string) error {
	_, err := r.db.Exec("DELETE FROM columns WHERE id = ?", id)
	return err
}

// Checklists
func (r *KanbanRepository) ListChecklists(taskID string) ([]models.Checklist, error) {
	var list []models.Checklist
	err := r.db.Select(&list, "SELECT id, task_id, title, position FROM checklists WHERE task_id = ? ORDER BY position ASC", taskID)
	if err != nil {
		return nil, err
	}
	for i := range list {
		var items []models.ChecklistItem
		err = r.db.Select(&items, "SELECT id, checklist_id, title, completed, position FROM checklist_items WHERE checklist_id = ? ORDER BY position ASC", list[i].ID)
		if err == nil {
			list[i].Items = items
		}
	}
	return list, nil
}

func (r *KanbanRepository) CreateChecklist(c *models.Checklist) error {
	_, err := r.db.Exec("INSERT INTO checklists (id, task_id, title, position) VALUES (?, ?, ?, ?)", c.ID, c.TaskID, c.Title, c.Position)
	return err
}

func (r *KanbanRepository) DeleteChecklist(id string) error {
	_, err := r.db.Exec("DELETE FROM checklists WHERE id = ?", id)
	return err
}

// ChecklistItems
func (r *KanbanRepository) CreateChecklistItem(item *models.ChecklistItem) error {
	_, err := r.db.Exec("INSERT INTO checklist_items (id, checklist_id, title, completed, position) VALUES (?, ?, ?, ?, ?)", item.ID, item.ChecklistID, item.Title, item.Completed, item.Position)
	return err
}

func (r *KanbanRepository) UpdateChecklistItem(item *models.ChecklistItem) error {
	_, err := r.db.Exec("UPDATE checklist_items SET title = ?, completed = ?, position = ? WHERE id = ?", item.Title, item.Completed, item.Position, item.ID)
	return err
}

func (r *KanbanRepository) DeleteChecklistItem(id string) error {
	_, err := r.db.Exec("DELETE FROM checklist_items WHERE id = ?", id)
	return err
}

// Labels
func (r *KanbanRepository) ListLabels() ([]models.Label, error) {
	var list []models.Label
	err := r.db.Select(&list, "SELECT id, name, color FROM task_labels")
	return list, err
}

func (r *KanbanRepository) SeedDefaultLabels() {
	labels := []models.Label{
		{ID: "w", Name: "Work", Color: "#3B82F6"},
		{ID: "s", Name: "Study", Color: "#10B981"},
		{ID: "p", Name: "Personal", Color: "#6366F1"},
		{ID: "u", Name: "Urgent", Color: "#EF4444"},
		{ID: "b", Name: "Bug", Color: "#EC4899"},
		{ID: "f", Name: "Feature", Color: "#F59E0B"},
		{ID: "r", Name: "Research", Color: "#8B5CF6"},
	}
	for _, l := range labels {
		_, _ = r.db.Exec("INSERT OR IGNORE INTO task_labels (id, name, color) VALUES (?, ?, ?)", l.ID, l.Name, l.Color)
	}
}

func (r *KanbanRepository) GetTaskLabels(taskID string) ([]models.Label, error) {
	var list []models.Label
	err := r.db.Select(&list, `
		SELECT l.id, l.name, l.color 
		FROM task_labels l
		JOIN task_label_relations r ON l.id = r.label_id
		WHERE r.task_id = ?
	`, taskID)
	return list, err
}

func (r *KanbanRepository) AddLabelToTask(taskID, labelID string) error {
	_, err := r.db.Exec("INSERT OR IGNORE INTO task_label_relations (task_id, label_id) VALUES (?, ?)", taskID, labelID)
	return err
}

func (r *KanbanRepository) RemoveLabelFromTask(taskID, labelID string) error {
	_, err := r.db.Exec("DELETE FROM task_label_relations WHERE task_id = ? AND label_id = ?", taskID, labelID)
	return err
}

// Attachments
func (r *KanbanRepository) ListAttachments(taskID string) ([]models.Attachment, error) {
	var list []models.Attachment
	err := r.db.Select(&list, "SELECT id, task_id, filename, url, size, created_at FROM attachments WHERE task_id = ? ORDER BY created_at DESC", taskID)
	return list, err
}

func (r *KanbanRepository) CreateAttachment(a *models.Attachment) error {
	_, err := r.db.Exec("INSERT INTO attachments (id, task_id, filename, url, size) VALUES (?, ?, ?, ?, ?)", a.ID, a.TaskID, a.Filename, a.URL, a.Size)
	return err
}

func (r *KanbanRepository) DeleteAttachment(id string) error {
	_, err := r.db.Exec("DELETE FROM attachments WHERE id = ?", id)
	return err
}

// Comments
func (r *KanbanRepository) ListComments(taskID string) ([]models.Comment, error) {
	var list []models.Comment
	err := r.db.Select(&list, "SELECT id, task_id, content, created_at, updated_at FROM comments WHERE task_id = ? ORDER BY created_at DESC", taskID)
	return list, err
}

func (r *KanbanRepository) CreateComment(c *models.Comment) error {
	_, err := r.db.Exec("INSERT INTO comments (id, task_id, content) VALUES (?, ?, ?)", c.ID, c.TaskID, c.Content)
	return err
}

func (r *KanbanRepository) DeleteComment(id string) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE id = ?", id)
	return err
}

// ActivityLogs
func (r *KanbanRepository) ListActivityLogs(taskID string) ([]models.ActivityLog, error) {
	var list []models.ActivityLog
	err := r.db.Select(&list, "SELECT id, task_id, type, old_value, new_value, created_at FROM activity_logs WHERE task_id = ? ORDER BY created_at DESC", taskID)
	return list, err
}

func (r *KanbanRepository) CreateActivityLog(log *models.ActivityLog) error {
	_, err := r.db.Exec("INSERT INTO activity_logs (id, task_id, type, old_value, new_value) VALUES (?, ?, ?, ?, ?)", log.ID, log.TaskID, log.Type, log.OldValue, log.NewValue)
	return err
}

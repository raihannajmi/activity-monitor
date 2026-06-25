package repositories

import (
	"fmt"
	"time"

	"activity-monitor/internal/models"

	"github.com/jmoiron/sqlx"
)

type NoteRepository struct {
	db *sqlx.DB
}

func NewNoteRepository(db *sqlx.DB) *NoteRepository {
	return &NoteRepository{db: db}
}

func (r *NoteRepository) List() ([]models.Note, error) {
	var notes []models.Note
	err := r.db.Select(&notes, `
		SELECT id, title, content, created_at, updated_at
		FROM notes ORDER BY updated_at DESC
	`)
	return notes, err
}

func (r *NoteRepository) Search(query string) ([]models.Note, error) {
	var notes []models.Note
	pattern := "%" + query + "%"
	err := r.db.Select(&notes, `
		SELECT id, title, content, created_at, updated_at
		FROM notes
		WHERE title LIKE ? OR content LIKE ?
		ORDER BY updated_at DESC
	`, pattern, pattern)
	return notes, err
}

func (r *NoteRepository) GetByID(id string) (*models.Note, error) {
	var note models.Note
	err := r.db.Get(&note, `
		SELECT id, title, content, created_at, updated_at FROM notes WHERE id=?
	`, id)
	if err != nil {
		return nil, fmt.Errorf("get note %s: %w", id, err)
	}
	return &note, nil
}

func (r *NoteRepository) Create(note *models.Note) error {
	_, err := r.db.Exec(`
		INSERT INTO notes (id, title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`, note.ID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt)
	return err
}

func (r *NoteRepository) Update(note *models.Note) error {
	note.UpdatedAt = time.Now()
	_, err := r.db.Exec(`
		UPDATE notes SET title=?, content=?, updated_at=? WHERE id=?
	`, note.Title, note.Content, note.UpdatedAt, note.ID)
	return err
}

func (r *NoteRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM notes WHERE id=?`, id)
	return err
}

func (r *NoteRepository) Count() (int, error) {
	var count int
	err := r.db.Get(&count, `SELECT COUNT(*) FROM notes`)
	return count, err
}

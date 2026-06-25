package repositories

import (
	"activity-monitor/internal/models"

	"github.com/jmoiron/sqlx"
)

type ActivityRepository struct {
	db *sqlx.DB
}

func NewActivityRepository(db *sqlx.DB) *ActivityRepository {
	return &ActivityRepository{db: db}
}

func (r *ActivityRepository) List(limit int) ([]models.Activity, error) {
	var activities []models.Activity
	err := r.db.Select(&activities, `
		SELECT id, type, reference_id, description, created_at
		FROM activities
		ORDER BY created_at DESC
		LIMIT ?
	`, limit)
	return activities, err
}

func (r *ActivityRepository) Create(a *models.Activity) error {
	_, err := r.db.Exec(`
		INSERT INTO activities (id, type, reference_id, description, created_at)
		VALUES (?, ?, ?, ?, ?)
	`, a.ID, a.Type, a.ReferenceID, a.Description, a.CreatedAt)
	return err
}

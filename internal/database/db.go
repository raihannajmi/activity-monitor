package database

import (
	"embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schemaFS embed.FS

func Open(path string) (*sqlx.DB, error) {
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if _, err := db.Exec("PRAGMA journal_mode=WAL; PRAGMA foreign_keys=ON;"); err != nil {
		return nil, fmt.Errorf("set pragmas: %w", err)
	}

	return db, nil
}

func Migrate(db *sqlx.DB) error {
	schema, err := schemaFS.ReadFile("schema.sql")
	if err != nil {
		return fmt.Errorf("read schema: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	// Safely add missing columns to tasks table for existing databases (ponytail: ignore errors if column already exists)
	alterCols := []string{
		"ALTER TABLE tasks ADD COLUMN column_id TEXT",
		"ALTER TABLE tasks ADD COLUMN position REAL DEFAULT 0.0",
		"ALTER TABLE tasks ADD COLUMN parent_task_id TEXT",
		"ALTER TABLE tasks ADD COLUMN estimate_pomodoro INTEGER DEFAULT 0",
		"ALTER TABLE tasks ADD COLUMN completed_pomodoro INTEGER DEFAULT 0",
		"ALTER TABLE tasks ADD COLUMN cover_image TEXT",
		"ALTER TABLE tasks ADD COLUMN created_by TEXT",
		"ALTER TABLE tasks ADD COLUMN updated_by TEXT",
		"ALTER TABLE tasks ADD COLUMN archived INTEGER DEFAULT 0",
	}
	for _, query := range alterCols {
		_, _ = db.Exec(query)
	}

	// Seed default board and columns if empty
	var count int
	if err := db.Get(&count, "SELECT COUNT(*) FROM boards"); err == nil && count == 0 {
		_, _ = db.Exec(`INSERT INTO boards (id, name, description, color) VALUES ('default', 'Main Board', 'Papan Kerja Utama', '#2563EB')`)
		_, _ = db.Exec(`INSERT INTO columns (id, board_id, name, position, color) VALUES ('todo', 'default', 'To Do', 1000.0, '#6B7280')`)
		_, _ = db.Exec(`INSERT INTO columns (id, board_id, name, position, color) VALUES ('inprogress', 'default', 'Sedang Dikerjakan', 2000.0, '#3B82F6')`)
		_, _ = db.Exec(`INSERT INTO columns (id, board_id, name, position, color) VALUES ('done', 'default', 'Selesai', 3000.0, '#10B981')`)
		
		// Map existing tasks
		_, _ = db.Exec(`UPDATE tasks SET column_id = 'todo' WHERE status = 'todo' AND column_id IS NULL`)
		_, _ = db.Exec(`UPDATE tasks SET column_id = 'inprogress' WHERE status = 'in_progress' AND column_id IS NULL`)
		_, _ = db.Exec(`UPDATE tasks SET column_id = 'done' WHERE status = 'done' AND column_id IS NULL`)
	}

	return nil
}

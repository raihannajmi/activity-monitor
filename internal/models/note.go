package models

import "time"

type Note struct {
	ID        string    `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (n *Note) Excerpt(max int) string {
	if len(n.Content) <= max {
		return n.Content
	}
	return n.Content[:max] + "..."
}

package services

import (
	"fmt"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"

	"github.com/google/uuid"
)

type NoteService struct {
	notes      *repositories.NoteRepository
	activities *repositories.ActivityRepository
}

func NewNoteService(notes *repositories.NoteRepository, activities *repositories.ActivityRepository) *NoteService {
	return &NoteService{notes, activities}
}

func (s *NoteService) List() ([]models.Note, error) {
	return s.notes.List()
}

func (s *NoteService) Search(query string) ([]models.Note, error) {
	if query == "" {
		return s.notes.List()
	}
	return s.notes.Search(query)
}

func (s *NoteService) GetByID(id string) (*models.Note, error) {
	return s.notes.GetByID(id)
}

func (s *NoteService) Create(title, content string) (*models.Note, error) {
	now := time.Now()
	note := &models.Note{
		ID:        uuid.New().String(),
		Title:     title,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.notes.Create(note); err != nil {
		return nil, fmt.Errorf("create note: %w", err)
	}
	s.logActivity(models.ActivityNoteCreated, note.ID, fmt.Sprintf(`Catatan "%s" dibuat`, note.Title))
	return note, nil
}

func (s *NoteService) Update(id, title, content string) (*models.Note, error) {
	note, err := s.notes.GetByID(id)
	if err != nil {
		return nil, err
	}
	note.Title = title
	note.Content = content
	if err := s.notes.Update(note); err != nil {
		return nil, err
	}
	s.logActivity(models.ActivityNoteUpdated, note.ID, fmt.Sprintf(`Catatan "%s" diperbarui`, note.Title))
	return note, nil
}

func (s *NoteService) Delete(id string) error {
	return s.notes.Delete(id)
}

func (s *NoteService) logActivity(actType models.ActivityType, refID, desc string) {
	a := &models.Activity{
		ID:          uuid.New().String(),
		Type:        actType,
		ReferenceID: refID,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	_ = s.activities.Create(a)
}

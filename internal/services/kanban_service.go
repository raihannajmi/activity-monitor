package services

import (
	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type KanbanService struct {
	repo       *repositories.KanbanRepository
	tasks      *repositories.TaskRepository
	activities *repositories.ActivityRepository
}

func NewKanbanService(repo *repositories.KanbanRepository, tasks *repositories.TaskRepository, activities *repositories.ActivityRepository) *KanbanService {
	repo.SeedDefaultLabels()
	return &KanbanService{repo: repo, tasks: tasks, activities: activities}
}

// Boards
func (s *KanbanService) ListBoards() ([]models.Board, error) {
	return s.repo.ListBoards()
}

func (s *KanbanService) CreateBoard(name, description, color string) (*models.Board, error) {
	b := &models.Board{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		Color:       color,
	}
	err := s.repo.CreateBoard(b)
	return b, err
}

// Columns
func (s *KanbanService) ListColumns(boardID string) ([]models.Column, error) {
	return s.repo.ListColumns(boardID)
}

func (s *KanbanService) CreateColumn(boardID, name string, position float64, color string) (*models.Column, error) {
	c := &models.Column{
		ID:       uuid.New().String(),
		BoardID:  boardID,
		Name:     name,
		Position: position,
		Color:    color,
	}
	err := s.repo.CreateColumn(c)
	return c, err
}

func (s *KanbanService) UpdateColumn(c *models.Column) error {
	return s.repo.UpdateColumn(c)
}

func (s *KanbanService) DeleteColumn(id string) error {
	return s.repo.DeleteColumn(id)
}

// Checklists
func (s *KanbanService) ListChecklists(taskID string) ([]models.Checklist, error) {
	return s.repo.ListChecklists(taskID)
}

func (s *KanbanService) CreateChecklist(taskID, title string, position float64) (*models.Checklist, error) {
	c := &models.Checklist{
		ID:       uuid.New().String(),
		TaskID:   taskID,
		Title:    title,
		Position: position,
	}
	err := s.repo.CreateChecklist(c)
	s.LogActivity(taskID, "create_checklist", "", title)
	return c, err
}

func (s *KanbanService) DeleteChecklist(taskID, id string) error {
	err := s.repo.DeleteChecklist(id)
	s.LogActivity(taskID, "delete_checklist", "", id)
	return err
}

// ChecklistItems
func (s *KanbanService) AddChecklistItem(checklistID, title string, position float64) (*models.ChecklistItem, error) {
	item := &models.ChecklistItem{
		ID:          uuid.New().String(),
		ChecklistID: checklistID,
		Title:       title,
		Completed:   false,
		Position:    position,
	}
	err := s.repo.CreateChecklistItem(item)
	return item, err
}

func (s *KanbanService) UpdateChecklistItem(item *models.ChecklistItem) error {
	return s.repo.UpdateChecklistItem(item)
}

func (s *KanbanService) DeleteChecklistItem(id string) error {
	return s.repo.DeleteChecklistItem(id)
}

// Labels
func (s *KanbanService) ListLabels() ([]models.Label, error) {
	return s.repo.ListLabels()
}

func (s *KanbanService) GetTaskLabels(taskID string) ([]models.Label, error) {
	return s.repo.GetTaskLabels(taskID)
}

func (s *KanbanService) AddLabelToTask(taskID, labelID string) error {
	return s.repo.AddLabelToTask(taskID, labelID)
}

func (s *KanbanService) RemoveLabelFromTask(taskID, labelID string) error {
	return s.repo.RemoveLabelFromTask(taskID, labelID)
}

// Attachments
func (s *KanbanService) ListAttachments(taskID string) ([]models.Attachment, error) {
	return s.repo.ListAttachments(taskID)
}

func (s *KanbanService) CreateAttachment(taskID, filename, url string, size int64) (*models.Attachment, error) {
	a := &models.Attachment{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		Filename:  filename,
		URL:       url,
		Size:      size,
		CreatedAt: time.Now(),
	}
	err := s.repo.CreateAttachment(a)
	s.LogActivity(taskID, "add_attachment", "", filename)
	return a, err
}

func (s *KanbanService) DeleteAttachment(taskID, id string) error {
	err := s.repo.DeleteAttachment(id)
	s.LogActivity(taskID, "delete_attachment", "", id)
	return err
}

// Comments
func (s *KanbanService) ListComments(taskID string) ([]models.Comment, error) {
	return s.repo.ListComments(taskID)
}

func (s *KanbanService) CreateComment(taskID, content string) (*models.Comment, error) {
	c := &models.Comment{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.repo.CreateComment(c)
	s.LogActivity(taskID, "add_comment", "", content)
	return c, err
}

func (s *KanbanService) DeleteComment(id string) error {
	return s.repo.DeleteComment(id)
}

// Activity Logging / Timeline
func (s *KanbanService) LogActivity(taskID string, activityType string, oldVal, newVal string) {
	log := &models.ActivityLog{
		ID:        uuid.New().String(),
		TaskID:    &taskID,
		Type:      activityType,
		OldValue:  &oldVal,
		NewValue:  &newVal,
		CreatedAt: time.Now(),
	}
	_ = s.repo.CreateActivityLog(log)

	// Also log into standard global timeline for synchronization
	desc := fmt.Sprintf("Task activity: %s -> %s", activityType, newVal)
	_ = s.activities.Create(&models.Activity{
		ID:          uuid.New().String(),
		Type:        models.ActivityType(activityType),
		ReferenceID: taskID,
		Description: desc,
		CreatedAt:   time.Now(),
	})
}

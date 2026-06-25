package services

import (
	"fmt"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"

	"github.com/google/uuid"
)

type ReminderService struct {
	reminders  *repositories.ReminderRepository
	tasks      *repositories.TaskRepository
	activities *repositories.ActivityRepository
}

func NewReminderService(
	reminders *repositories.ReminderRepository,
	tasks *repositories.TaskRepository,
	activities *repositories.ActivityRepository,
) *ReminderService {
	return &ReminderService{reminders, tasks, activities}
}

func (s *ReminderService) ListToday() ([]models.Reminder, error) {
	reminders, err := s.reminders.ListToday()
	if err != nil {
		return nil, err
	}
	s.populateTaskTitles(reminders)
	return reminders, nil
}

func (s *ReminderService) ListUpcoming() ([]models.Reminder, error) {
	reminders, err := s.reminders.ListUpcoming()
	if err != nil {
		return nil, err
	}
	s.populateTaskTitles(reminders)
	return reminders, nil
}

func (s *ReminderService) Create(taskID, note string, remindAt time.Time) (*models.Reminder, error) {
	rem := &models.Reminder{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		RemindAt:  remindAt,
		Note:      note,
		CreatedAt: time.Now(),
	}
	if err := s.reminders.Create(rem); err != nil {
		return nil, fmt.Errorf("create reminder: %w", err)
	}
	s.logActivity(models.ActivityReminderCreated, rem.ID, fmt.Sprintf(`Reminder dibuat untuk %s`, remindAt.Format("02 Jan 2006 15:04")))
	return rem, nil
}

func (s *ReminderService) MarkDone(id string) error {
	return s.reminders.MarkDone(id)
}

func (s *ReminderService) Delete(id string) error {
	return s.reminders.Delete(id)
}

func (s *ReminderService) CountToday() (int, error) {
	return s.reminders.CountToday()
}

func (s *ReminderService) populateTaskTitles(reminders []models.Reminder) {
	for i := range reminders {
		if reminders[i].TaskID != "" {
			if task, err := s.tasks.GetByID(reminders[i].TaskID); err == nil {
				reminders[i].TaskTitle = task.Title
			}
		}
	}
}

func (s *ReminderService) logActivity(actType models.ActivityType, refID, desc string) {
	a := &models.Activity{
		ID:          uuid.New().String(),
		Type:        actType,
		ReferenceID: refID,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	_ = s.activities.Create(a)
}

package services

import (
	"fmt"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"

	"github.com/google/uuid"
)

type TaskService struct {
	tasks      *repositories.TaskRepository
	subtasks   *repositories.SubtaskRepository
	reminders  *repositories.ReminderRepository
	activities *repositories.ActivityRepository
	timelogs   *repositories.TimeLogRepository
}

func NewTaskService(
	tasks *repositories.TaskRepository,
	subtasks *repositories.SubtaskRepository,
	reminders *repositories.ReminderRepository,
	activities *repositories.ActivityRepository,
	timelogs *repositories.TimeLogRepository,
) *TaskService {
	return &TaskService{tasks, subtasks, reminders, activities, timelogs}
}

func (s *TaskService) ListByStatus(status string) ([]models.Task, error) {
	return s.tasks.ListByStatus(status)
}

func (s *TaskService) GetSubtasks(taskID string) ([]models.Subtask, error) {
	subs, err := s.subtasks.ListByTaskID(taskID)
	if err != nil {
		return nil, err
	}
	for i := range subs {
		subs[i].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForSubtask(subs[i].ID)
	}
	return subs, nil
}

func (s *TaskService) ListAll() ([]models.Task, error) {
	return s.tasks.ListAll()
}

func (s *TaskService) ListWithSubtasks() ([]models.Task, error) {
	tasks, err := s.tasks.ListAll()
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		subs, _ := s.subtasks.ListByTaskID(tasks[i].ID)
		for j := range subs {
			subs[j].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForSubtask(subs[j].ID)
		}
		tasks[i].Subtasks = subs
		tasks[i].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForTask(tasks[i].ID)
	}
	return tasks, nil
}

func (s *TaskService) GetWithDetails(id string) (*models.Task, error) {
	task, err := s.tasks.GetByID(id)
	if err != nil {
		return nil, err
	}
	subs, _ := s.subtasks.ListByTaskID(id)
	for j := range subs {
		subs[j].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForSubtask(subs[j].ID)
	}
	task.Subtasks = subs
	task.Reminders, _ = s.reminders.ListByTaskID(id)
	task.TotalTimeSpent, _ = s.timelogs.GetTotalTimeForTask(id)
	return task, nil
}

func (s *TaskService) ListDueTodayWithSubtasks() ([]models.Task, error) {
	tasks, err := s.tasks.ListDueToday()
	if err != nil {
		return nil, err
	}
	for i := range tasks {
		subs, _ := s.subtasks.ListByTaskID(tasks[i].ID)
		for j := range subs {
			subs[j].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForSubtask(subs[j].ID)
		}
		tasks[i].Subtasks = subs
		tasks[i].TotalTimeSpent, _ = s.timelogs.GetTotalTimeForTask(tasks[i].ID)
	}
	return tasks, nil
}

func (s *TaskService) Create(title, description string, priority models.Priority, deadline *time.Time) (*models.Task, error) {
	task := &models.Task{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      models.StatusTodo,
		Deadline:    deadline,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.tasks.Create(task); err != nil {
		return nil, fmt.Errorf("create task: %w", err)
	}
	s.logActivity(models.ActivityTaskCreated, task.ID, fmt.Sprintf(`Task "%s" dibuat`, task.Title))
	return task, nil
}

func (s *TaskService) Update(task *models.Task) error {
	if err := s.tasks.Update(task); err != nil {
		return err
	}
	s.logActivity(models.ActivityTaskUpdated, task.ID, fmt.Sprintf(`Task "%s" diperbarui`, task.Title))
	return nil
}

func (s *TaskService) Delete(id string) error {
	task, err := s.tasks.GetByID(id)
	if err != nil {
		return err
	}
	if err := s.tasks.Delete(id); err != nil {
		return err
	}
	s.logActivity(models.ActivityTaskUpdated, id, fmt.Sprintf(`Task "%s" dihapus`, task.Title))
	return nil
}

func (s *TaskService) AddSubtask(taskID, title string) (*models.Subtask, error) {
	sub := &models.Subtask{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		Title:     title,
		CreatedAt: time.Now(),
	}
	if err := s.subtasks.Create(sub); err != nil {
		return nil, err
	}
	return sub, nil
}

func (s *TaskService) ToggleSubtask(subtaskID, taskID string) error {
	_, err := s.subtasks.ToggleComplete(subtaskID)
	if err != nil {
		return err
	}

	// Check if all subtasks are done → auto-complete task
	allDone, err := s.subtasks.AllCompleted(taskID)
	if err != nil {
		return err
	}
	if allDone {
		task, err := s.tasks.GetByID(taskID)
		if err != nil {
			return err
		}
		if task.Status != models.StatusDone {
			if err := s.tasks.UpdateStatus(taskID, models.StatusDone); err != nil {
				return err
			}
			s.logActivity(models.ActivityTaskDone, taskID, fmt.Sprintf(`Task "%s" selesai`, task.Title))
		}
	}
	return nil
}

func (s *TaskService) DeleteSubtask(id string) error {
	return s.subtasks.Delete(id)
}

func (s *TaskService) UpdateStatus(id string, status models.Status) error {
	task, err := s.tasks.GetByID(id)
	if err != nil {
		return err
	}
	if err := s.tasks.UpdateStatus(id, status); err != nil {
		return err
	}
	if status == models.StatusDone {
		s.logActivity(models.ActivityTaskDone, id, fmt.Sprintf(`Task "%s" selesai`, task.Title))
	} else {
		s.logActivity(models.ActivityTaskUpdated, id, fmt.Sprintf(`Task "%s" diperbarui`, task.Title))
	}
	return nil
}

func (s *TaskService) GetStats() (active, done, today int, err error) {
	active, err = s.tasks.CountActive()
	if err != nil {
		return
	}
	done, err = s.tasks.CountByStatus(string(models.StatusDone))
	if err != nil {
		return
	}
	todayTasks, err := s.tasks.ListDueToday()
	today = len(todayTasks)
	return
}

func (s *TaskService) logActivity(actType models.ActivityType, refID, desc string) {
	a := &models.Activity{
		ID:          uuid.New().String(),
		Type:        actType,
		ReferenceID: refID,
		Description: desc,
		CreatedAt:   time.Now(),
	}
	_ = s.activities.Create(a)
}

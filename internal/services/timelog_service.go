package services

import (
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"
)

type TimeLogService struct {
	timelogs *repositories.TimeLogRepository
}

func NewTimeLogService(timelogs *repositories.TimeLogRepository) *TimeLogService {
	return &TimeLogService{timelogs: timelogs}
}

func (s *TimeLogService) StartSession(taskID, subtaskID *string, sessionType string) (*models.TimeLog, error) {
	// If there's an active session, stop it first?
	// For simplicity, we just let the repo start a new one. Or we could stop the existing one.
	active, err := s.timelogs.GetActiveSession()
	if err == nil && active != nil {
		duration := int(time.Since(active.StartTime).Seconds())
		s.timelogs.StopSession(active.ID, duration)
	}

	return s.timelogs.StartSession(taskID, subtaskID, sessionType)
}

func (s *TimeLogService) StopSession(id string, durationSeconds int) error {
	return s.timelogs.StopSession(id, durationSeconds)
}

func (s *TimeLogService) GetActiveSession() (*models.TimeLog, error) {
	return s.timelogs.GetActiveSession()
}

func (s *TimeLogService) GetDailyRecap(date time.Time) (*models.DailyRecap, error) {
	return s.timelogs.GetDailyRecap(date)
}

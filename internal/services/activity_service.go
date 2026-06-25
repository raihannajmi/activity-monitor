package services

import (
	"activity-monitor/internal/models"
	"activity-monitor/internal/repositories"
)

type ActivityService struct {
	activities *repositories.ActivityRepository
}

func NewActivityService(activities *repositories.ActivityRepository) *ActivityService {
	return &ActivityService{activities}
}

func (s *ActivityService) List(limit int) ([]models.Activity, error) {
	if limit <= 0 {
		limit = 50
	}
	return s.activities.List(limit)
}

func (s *ActivityService) ListRecent() ([]models.Activity, error) {
	return s.activities.List(10)
}

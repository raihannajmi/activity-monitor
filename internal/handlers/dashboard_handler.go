package handlers

import (
	"net/http"
	"time"

	"activity-monitor/internal/services"
	"activity-monitor/templates/pages"
)

type DashboardHandler struct {
	tasks      *services.TaskService
	reminders  *services.ReminderService
	activities *services.ActivityService
	timelogs   *services.TimeLogService
}

func NewDashboardHandler(tasks *services.TaskService, reminders *services.ReminderService, activities *services.ActivityService, timelogs *services.TimeLogService) *DashboardHandler {
	return &DashboardHandler{tasks, reminders, activities, timelogs}
}

func (h *DashboardHandler) Show(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	active, done, today, err := h.tasks.GetStats()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todayReminders, _ := h.reminders.CountToday()
	tasksDueToday, _ := h.tasks.ListDueTodayWithSubtasks()
	reminders, _ := h.reminders.ListToday()
	
	// Get daily recap instead of recent activity
	recap, _ := h.timelogs.GetDailyRecap(time.Now())

	data := pages.DashboardData{
		ActiveTasks:    active,
		DoneTasks:      done,
		TodayTasks:     today,
		TodayReminders: todayReminders,
		TasksDueToday:  tasksDueToday,
		Reminders:      reminders,
		DailyRecap:     recap,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Dashboard(data).Render(r.Context(), w)
}

package handlers

import (
	"net/http"

	"activity-monitor/internal/services"
	"activity-monitor/templates/pages"
)

type DashboardHandler struct {
	tasks      *services.TaskService
	reminders  *services.ReminderService
	activities *services.ActivityService
}

func NewDashboardHandler(tasks *services.TaskService, reminders *services.ReminderService, activities *services.ActivityService) *DashboardHandler {
	return &DashboardHandler{tasks, reminders, activities}
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
	recent, _ := h.activities.ListRecent()

	data := pages.DashboardData{
		ActiveTasks:    active,
		DoneTasks:      done,
		TodayTasks:     today,
		TodayReminders: todayReminders,
		TasksDueToday:  tasksDueToday,
		Reminders:      reminders,
		RecentActivity: recent,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Dashboard(data).Render(r.Context(), w)
}

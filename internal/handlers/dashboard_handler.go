package handlers

import (
	"net/http"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/services"
	"activity-monitor/templates/components"
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

	allTasks, _ := h.tasks.ListWithSubtasks()
	allLogs, _ := h.timelogs.ListAllWithDetails()
	
	// Get top 5 recent logs
	var recentLogs []models.TimeLogWithDetails
	if len(allLogs) > 5 {
		recentLogs = allLogs[:5]
	} else {
		recentLogs = allLogs
	}

	data := pages.DashboardData{
		ActiveTasks:    active,
		DoneTasks:      done,
		TodayTasks:     today,
		TodayReminders: todayReminders,
		TasksDueToday:  tasksDueToday,
		Reminders:      reminders,
		DailyRecap:     recap,
		AllTasks:       allTasks,
		RecentTimeLogs: recentLogs,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Dashboard(data).Render(r.Context(), w)
}

func (h *DashboardHandler) SidebarStats(w http.ResponseWriter, r *http.Request) {
	recap, _ := h.timelogs.GetDailyRecap(time.Now())
	
	// Assuming target is 5 sessions and 2 hours (7200 seconds)
	targetSessions := 5
	targetSeconds := 7200

	progress := 0
	if recap != nil && recap.TotalFocusSeconds > 0 {
		progress = int(float64(recap.TotalFocusSeconds) / float64(targetSeconds) * 100)
		if progress > 100 {
			progress = 100
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.SidebarFocusWidget(recap, targetSessions, targetSeconds, progress).Render(r.Context(), w)
}

package handlers

import (
	"net/http"
	"strings"

	"activity-monitor/internal/services"
	"activity-monitor/templates/pages"
)

type TimelineHandler struct {
	activities *services.ActivityService
}

func NewTimelineHandler(activities *services.ActivityService) *TimelineHandler {
	return &TimelineHandler{activities}
}

func (h *TimelineHandler) Show(w http.ResponseWriter, r *http.Request) {
	activities, err := h.activities.List(100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Timeline(activities).Render(r.Context(), w)
}

type ReminderHandler struct {
	reminders *services.ReminderService
}

func NewReminderHandler(reminders *services.ReminderService) *ReminderHandler {
	return &ReminderHandler{reminders}
}

func (h *ReminderHandler) MarkDone(w http.ResponseWriter, r *http.Request) {
	// Path: /reminders/{id}/done
	path := strings.TrimPrefix(r.URL.Path, "/reminders/")
	id := strings.TrimSuffix(path, "/done")

	if err := h.reminders.MarkDone(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return empty to remove from view
	w.WriteHeader(http.StatusOK)
}

func (h *ReminderHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/reminders/")
	if err := h.reminders.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

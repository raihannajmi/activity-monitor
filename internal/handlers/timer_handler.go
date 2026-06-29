package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"activity-monitor/internal/services"
)

type TimerHandler struct {
	timelogs *services.TimeLogService
}

func NewTimerHandler(timelogs *services.TimeLogService) *TimerHandler {
	return &TimerHandler{timelogs: timelogs}
}

func (h *TimerHandler) Start(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	taskID := r.FormValue("task_id")
	subtaskID := r.FormValue("subtask_id")
	sessionType := r.FormValue("session_type")
	if sessionType == "" {
		sessionType = "pomodoro"
	}

	var pTaskID, pSubtaskID *string
	if taskID != "" {
		pTaskID = &taskID
	}
	if subtaskID != "" {
		pSubtaskID = &subtaskID
	}

	session, err := h.timelogs.StartSession(pTaskID, pSubtaskID, sessionType)
	if err != nil {
		http.Error(w, "Gagal memulai timer", http.StatusInternalServerError)
		return
	}

	// Emit an HTMX event to trigger the timer widget
	w.Header().Set("HX-Trigger", fmt.Sprintf(`{"timerStarted": {"id": "%s", "type": "%s", "task_id": "%s", "subtask_id": "%s"}}`, session.ID, session.SessionType, taskID, subtaskID))
	w.WriteHeader(http.StatusOK)
}

func (h *TimerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/timer/")
	id = strings.TrimSuffix(id, "/stop")

	durationStr := r.FormValue("duration_seconds")
	duration, _ := strconv.Atoi(durationStr)

	if err := h.timelogs.StopSession(id, duration); err != nil {
		http.Error(w, "Gagal menghentikan timer", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Trigger", `{"timerStopped": true, "statsUpdated": true}`)
	w.WriteHeader(http.StatusOK)
}

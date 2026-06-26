package handlers

import (
	"net/http"

	"activity-monitor/internal/services"
	"activity-monitor/templates/pages"
)

type ReportHandler struct {
	timelogs *services.TimeLogService
}

func NewReportHandler(timelogs *services.TimeLogService) *ReportHandler {
	return &ReportHandler{timelogs: timelogs}
}

func (h *ReportHandler) Show(w http.ResponseWriter, r *http.Request) {
	logs, err := h.timelogs.ListAllWithDetails()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Reports(logs).Render(r.Context(), w)
}

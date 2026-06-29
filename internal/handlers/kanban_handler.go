package handlers

import (
	"activity-monitor/internal/models"
	"activity-monitor/internal/services"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type KanbanHandler struct {
	kanban *services.KanbanService
	tasks  *services.TaskService
}

func NewKanbanHandler(kanban *services.KanbanService, tasks *services.TaskService) *KanbanHandler {
	return &KanbanHandler{kanban: kanban, tasks: tasks}
}

func (h *KanbanHandler) MoveTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/move")

	task, err := h.tasks.GetWithDetails(taskID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	columnID := r.FormValue("column_id")
	prevID := r.FormValue("prev_id")
	nextID := r.FormValue("next_id")

	// Calculate floating position (Trello-style position arithmetic)
	var newPos float64
	var posPrev, posNext float64
	hasPrev, hasNext := false, false

	if prevID != "" {
		if prevTask, err := h.tasks.GetWithDetails(prevID); err == nil {
			posPrev = prevTask.Position
			hasPrev = true
		}
	}
	if nextID != "" {
		if nextTask, err := h.tasks.GetWithDetails(nextID); err == nil {
			posNext = nextTask.Position
			hasNext = true
		}
	}

	if hasPrev && hasNext {
		newPos = (posPrev + posNext) / 2.0
	} else if hasPrev {
		newPos = posPrev + 1000.0
	} else if hasNext {
		newPos = posNext / 2.0
	} else {
		newPos = 1000.0
	}

	task.Position = newPos
	if columnID != "" {
		task.ColumnID = &columnID
		// Sync task status based on column ID
		if columnID == "todo" {
			task.Status = models.StatusTodo
		} else if columnID == "inprogress" || columnID == "in_progress" {
			task.Status = models.StatusInProgress
		} else if columnID == "done" {
			task.Status = models.StatusDone
		}
	}

	if err := h.tasks.Update(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.kanban.LogActivity(taskID, "move_task", "", fmt.Sprintf("Moved to column %s at position %f", columnID, newPos))

	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/comments")

	content := strings.TrimSpace(r.FormValue("content"))
	if content == "" {
		http.Error(w, "Comment content required", http.StatusBadRequest)
		return
	}

	_, err := h.kanban.CreateComment(taskID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trigger WebSocket/HTMX reload of comments section
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) AddChecklist(w http.ResponseWriter, r *http.Request) {
	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/checklists")

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		title = "Checklist"
	}

	_, err := h.kanban.CreateChecklist(taskID, title, 1000.0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) AddAttachment(w http.ResponseWriter, r *http.Request) {
	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/attachments")

	filename := r.FormValue("filename")
	url := r.FormValue("url")
	if filename == "" || url == "" {
		http.Error(w, "Filename and URL required", http.StatusBadRequest)
		return
	}

	_, err := h.kanban.CreateAttachment(taskID, filename, url, 100)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) ToggleLabel(w http.ResponseWriter, r *http.Request) {
	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/labels")

	labelID := r.FormValue("label_id")
	if labelID == "" {
		http.Error(w, "Label ID required", http.StatusBadRequest)
		return
	}

	// Simple toggle: remove if exists, add if not
	labels, _ := h.kanban.GetTaskLabels(taskID)
	exists := false
	for _, l := range labels {
		if l.ID == labelID {
			exists = true
			break
		}
	}

	var err error
	if exists {
		err = h.kanban.RemoveLabelFromTask(taskID, labelID)
	} else {
		err = h.kanban.AddLabelToTask(taskID, labelID)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

func (h *KanbanHandler) UpdatePomodoro(w http.ResponseWriter, r *http.Request) {
	taskID := extractID(r.URL.Path, "/tasks/")
	taskID = strings.TrimSuffix(taskID, "/pomodoro")

	task, err := h.tasks.GetWithDetails(taskID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	est, _ := strconv.Atoi(r.FormValue("estimate_pomodoro"))
	comp, _ := strconv.Atoi(r.FormValue("completed_pomodoro"))

	task.EstimatePomodoro = est
	task.CompletedPomodoro = comp

	_ = h.tasks.Update(task)
	w.Header().Set("HX-Refresh", "true")
	w.WriteHeader(http.StatusOK)
}

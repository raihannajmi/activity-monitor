package handlers

import (
	"net/http"
	"strings"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/services"
	"activity-monitor/templates/components"
	"activity-monitor/templates/pages"
)

type TaskHandler struct {
	tasks     *services.TaskService
	reminders *services.ReminderService
}

func NewTaskHandler(tasks *services.TaskService, reminders *services.ReminderService) *TaskHandler {
	return &TaskHandler{tasks, reminders}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	var taskList []models.Task
	var err error

	if filter != "" && filter != "all" {
		taskList, err = h.tasks.ListByStatus(filter)
	} else {
		taskList, err = h.tasks.ListWithSubtasks()
	}
	if err != nil {
		h.sendError(w, r, "Gagal mengambil data task", http.StatusInternalServerError)
		return
	}

	// Populate subtasks for filtered results
	if filter != "" && filter != "all" {
		for i := range taskList {
			subtasks, _ := h.tasks.GetSubtasks(taskList[i].ID)
			taskList[i].Subtasks = subtasks
		}
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Tasks(taskList, filter).Render(r.Context(), w)
}

func (h *TaskHandler) Detail(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/tasks/")
	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.TaskDetail(*task).Render(r.Context(), w)
}

func (h *TaskHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.TaskFormModal(nil, nil).Render(r.Context(), w)
}

func (h *TaskHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/tasks/")
	id = strings.TrimSuffix(id, "/edit")
	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.TaskFormModal(task, nil).Render(r.Context(), w)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data tidak valid", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		h.sendError(w, r, "Judul task tidak boleh kosong", http.StatusBadRequest)
		return
	}

	priority := models.Priority(r.FormValue("priority"))
	if priority == "" {
		priority = models.PriorityMedium
	}

	var deadline *time.Time
	if d := r.FormValue("deadline"); d != "" {
		if t, err := time.ParseInLocation("2006-01-02", d, time.Local); err == nil {
			deadline = &t
		}
	}

	if _, err := h.tasks.Create(title, r.FormValue("description"), priority, deadline); err != nil {
		h.sendError(w, r, "Gagal membuat task", http.StatusInternalServerError)
		return
	}

	// Re-render full task list
	taskList, _ := h.tasks.ListWithSubtasks()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTaskList(w, r, taskList)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/tasks/")
	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data form tidak valid", http.StatusBadRequest)
		return
	}

	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	task.Title = strings.TrimSpace(r.FormValue("title"))
	task.Description = r.FormValue("description")
	task.Priority = models.Priority(r.FormValue("priority"))

	var deadline *time.Time
	if d := r.FormValue("deadline"); d != "" {
		if t, err := time.ParseInLocation("2006-01-02", d, time.Local); err == nil {
			deadline = &t
		}
	}
	task.Deadline = deadline

	if err := h.tasks.Update(task); err != nil {
		h.sendError(w, r, "Gagal menyimpan task", http.StatusInternalServerError)
		return
	}

	// Context-aware update: if edited from Detail Page, refresh to see changes properly
	if strings.Contains(r.Header.Get("HX-Current-URL"), "/tasks/"+id) {
		w.Header().Set("HX-Refresh", "true")
		return
	}

	taskList, _ := h.tasks.ListWithSubtasks()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	renderTaskList(w, r, taskList)
}

func (h *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{id}/status
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		h.sendError(w, r, "Path tidak valid", http.StatusBadRequest)
		return
	}
	id := parts[0]

	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data tidak valid", http.StatusBadRequest)
		return
	}

	status := models.Status(r.FormValue("status"))
	if err := h.tasks.UpdateStatus(id, status); err != nil {
		h.sendError(w, r, "Gagal mengubah status task", http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.TaskCard(*task).Render(r.Context(), w)
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/tasks/")
	if err := h.tasks.Delete(id); err != nil {
		h.sendError(w, r, "Gagal menghapus task", http.StatusInternalServerError)
		return
	}
	
	// If deleted from Detail Page, redirect back to task list
	if strings.Contains(r.Header.Get("HX-Current-URL"), "/tasks/"+id) {
		w.Header().Set("HX-Redirect", "/tasks")
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) AddSubtask(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{id}/subtasks
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	taskID := parts[0]

	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data tidak valid", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		h.sendError(w, r, "Judul subtask tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if _, err := h.tasks.AddSubtask(taskID, title); err != nil {
		h.sendError(w, r, "Gagal menambahkan subtask", http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetWithDetails(taskID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.SubtaskList(*task).Render(r.Context(), w)
}

func (h *TaskHandler) ToggleSubtask(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{taskID}/subtasks/{subtaskID}/toggle
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		h.sendError(w, r, "Path tidak valid", http.StatusBadRequest)
		return
	}
	taskID := parts[0]
	subtaskID := parts[2]

	if err := h.tasks.ToggleSubtask(subtaskID, taskID); err != nil {
		h.sendError(w, r, "Gagal mengubah status subtask", http.StatusInternalServerError)
		return
	}

	task, err := h.tasks.GetWithDetails(taskID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.SubtaskList(*task).Render(r.Context(), w)
}

func (h *TaskHandler) DeleteSubtask(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{taskID}/subtasks/{subtaskID}
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		h.sendError(w, r, "Path tidak valid", http.StatusBadRequest)
		return
	}
	subtaskID := parts[2]

	if err := h.tasks.DeleteSubtask(subtaskID); err != nil {
		h.sendError(w, r, "Gagal menghapus subtask", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *TaskHandler) ReminderForm(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{id}/reminder/new
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	taskID := parts[0]

	task, err := h.tasks.GetWithDetails(taskID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.ReminderFormModal(taskID, task.Title).Render(r.Context(), w)
}

func (h *TaskHandler) CreateReminder(w http.ResponseWriter, r *http.Request) {
	// Path: /tasks/{id}/reminders
	path := strings.TrimPrefix(r.URL.Path, "/tasks/")
	parts := strings.Split(path, "/")
	taskID := parts[0]

	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data form tidak valid", http.StatusBadRequest)
		return
	}

	dateStr := r.FormValue("date")
	timeStr := r.FormValue("time")
	note := r.FormValue("note")

	remindAt, err := time.ParseInLocation("2006-01-02 15:04", dateStr+" "+timeStr, time.Local)
	if err != nil {
		h.sendError(w, r, "Format tanggal atau jam tidak valid", http.StatusBadRequest)
		return
	}

	if _, err := h.reminders.Create(taskID, note, remindAt); err != nil {
		h.sendError(w, r, "Gagal membuat reminder", http.StatusInternalServerError)
		return
	}

	// Redirect back to task detail
	http.Redirect(w, r, "/tasks/"+taskID, http.StatusSeeOther)
}

// Helper: render task list div for HTMX swaps
func renderTaskList(w http.ResponseWriter, r *http.Request, tasks []models.Task) {
	w.Write([]byte(`<div class="task-list" id="task-list">`))
	for _, t := range tasks {
		components.TaskCard(t).Render(r.Context(), w)
	}
	w.Write([]byte(`</div>`))
}

func extractID(path, prefix string) string {
	id := strings.TrimPrefix(path, prefix)
	// Remove any trailing path segments
	if idx := strings.Index(id, "/"); idx != -1 {
		id = id[:idx]
	}
	return id
}

func (h *TaskHandler) sendError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	w.Header().Set("HX-Retarget", "#toast-container")
	w.Header().Set("HX-Reswap", "beforeend")
	w.WriteHeader(http.StatusOK)
	components.ToastError(msg).Render(r.Context(), w)
}

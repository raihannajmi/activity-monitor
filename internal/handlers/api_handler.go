package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"activity-monitor/internal/models"
	"activity-monitor/internal/services"
)

type APIHandler struct {
	tasks      *services.TaskService
	notes      *services.NoteService
	reminders  *services.ReminderService
	activities *services.ActivityService
	timelogs   *services.TimeLogService
	kanban     *services.KanbanService
}

func NewAPIHandler(
	tasks *services.TaskService,
	notes *services.NoteService,
	reminders *services.ReminderService,
	activities *services.ActivityService,
	timelogs *services.TimeLogService,
	kanban *services.KanbanService,
) *APIHandler {
	return &APIHandler{
		tasks:      tasks,
		notes:      notes,
		reminders:  reminders,
		activities: activities,
		timelogs:   timelogs,
		kanban:     kanban,
	}
}

// Helpers
func (h *APIHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func (h *APIHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

func (h *APIHandler) extractID(path, prefix string) string {
	id := strings.TrimPrefix(path, prefix)
	if idx := strings.Index(id, "/"); idx != -1 {
		id = id[:idx]
	}
	return id
}

// Dashboard
func (h *APIHandler) DashboardStats(w http.ResponseWriter, r *http.Request) {
	active, done, today, err := h.tasks.GetStats()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	todayReminders, _ := h.reminders.CountToday()
	tasksDueToday, _ := h.tasks.ListDueTodayWithSubtasks()
	remindersList, _ := h.reminders.ListToday()
	recap, _ := h.timelogs.GetDailyRecap(time.Now())
	allLogs, _ := h.timelogs.ListAllWithDetails()

	var recentLogs []models.TimeLogWithDetails
	if len(allLogs) > 5 {
		recentLogs = allLogs[:5]
	} else {
		recentLogs = allLogs
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"active_tasks":    active,
		"done_tasks":      done,
		"today_tasks":     today,
		"today_reminders": todayReminders,
		"tasks_due_today": tasksDueToday,
		"reminders":       remindersList,
		"daily_recap":     recap,
		"recent_logs":     recentLogs,
	})
}

// Tasks
func (h *APIHandler) TaskList(w http.ResponseWriter, r *http.Request) {
	filter := r.URL.Query().Get("filter")
	var taskList []models.Task
	var err error

	if filter != "" && filter != "all" {
		taskList, err = h.tasks.ListByStatus(filter)
	} else {
		taskList, err = h.tasks.ListWithSubtasks()
	}
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if filter != "" && filter != "all" {
		for i := range taskList {
			subtasks, _ := h.tasks.GetSubtasks(taskList[i].ID)
			taskList[i].Subtasks = subtasks
		}
	}
	h.writeJSON(w, http.StatusOK, taskList)
}

func (h *APIHandler) TaskCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    string `json:"priority"`
		Deadline    string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		h.writeError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	prio := models.PriorityMedium
	if input.Priority != "" {
		prio = models.Priority(input.Priority)
	}

	var deadline *time.Time
	if input.Deadline != "" {
		if t, err := time.Parse("2006-01-02", input.Deadline); err == nil {
			deadline = &t
		}
	}

	task, err := h.tasks.Create(title, input.Description, prio, deadline)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, task)
}

func (h *APIHandler) TaskGet(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/tasks/")
	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found")
		return
	}
	h.writeJSON(w, http.StatusOK, task)
}

func (h *APIHandler) TaskUpdate(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/tasks/")
	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	var input struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		Priority          string `json:"priority"`
		Deadline          string `json:"deadline"`
		Status            string `json:"status"`
		EstimatePomodoro  int    `json:"estimate_pomodoro"`
		CompletedPomodoro int    `json:"completed_pomodoro"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	if strings.TrimSpace(input.Title) != "" {
		task.Title = strings.TrimSpace(input.Title)
	}
	task.Description = input.Description
	if input.Priority != "" {
		task.Priority = models.Priority(input.Priority)
	}
	if input.Status != "" {
		task.Status = models.Status(input.Status)
	}

	var deadline *time.Time
	if input.Deadline != "" {
		if t, err := time.Parse("2006-01-02", input.Deadline); err == nil {
			deadline = &t
		}
	}
	task.Deadline = deadline
	task.EstimatePomodoro = input.EstimatePomodoro
	task.CompletedPomodoro = input.CompletedPomodoro

	if err := h.tasks.Update(task); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, task)
}

func (h *APIHandler) TaskDelete(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/tasks/")
	if err := h.tasks.Delete(id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *APIHandler) TaskMove(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/tasks/")
	task, err := h.tasks.GetWithDetails(id)
	if err != nil {
		h.writeError(w, http.StatusNotFound, "Task not found")
		return
	}

	var input struct {
		ColumnID string `json:"column_id"`
		PrevID   string `json:"prev_id"`
		NextID   string `json:"next_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	var newPos float64
	var posPrev, posNext float64
	hasPrev, hasNext := false, false

	if input.PrevID != "" {
		if prevTask, err := h.tasks.GetWithDetails(input.PrevID); err == nil {
			posPrev = prevTask.Position
			hasPrev = true
		}
	}
	if input.NextID != "" {
		if nextTask, err := h.tasks.GetWithDetails(input.NextID); err == nil {
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
	if input.ColumnID != "" {
		task.ColumnID = &input.ColumnID
		if input.ColumnID == "todo" {
			task.Status = models.StatusTodo
		} else if input.ColumnID == "inprogress" || input.ColumnID == "in_progress" {
			task.Status = models.StatusInProgress
		} else if input.ColumnID == "done" {
			task.Status = models.StatusDone
		}
	}

	if err := h.tasks.Update(task); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.kanban.LogActivity(id, "move_task", "", fmt.Sprintf("Moved to column %s", input.ColumnID))
	h.writeJSON(w, http.StatusOK, task)
}

// Subtasks
func (h *APIHandler) SubtaskCreate(w http.ResponseWriter, r *http.Request) {
	taskID := h.extractID(r.URL.Path, "/api/v1/tasks/")
	var input struct {
		Title    string `json:"title"`
		Deadline string `json:"deadline"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		h.writeError(w, http.StatusBadRequest, "Subtask title cannot be empty")
		return
	}

	var deadline *time.Time
	if input.Deadline != "" {
		if t, err := time.Parse("2006-01-02", input.Deadline); err == nil {
			deadline = &t
		}
	}

	sub, err := h.tasks.AddSubtask(taskID, title, deadline)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, sub)
}

func (h *APIHandler) SubtaskToggle(w http.ResponseWriter, r *http.Request) {
	// Path: /api/v1/tasks/{task_id}/subtasks/{subtask_id}/toggle
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	taskID := parts[4]
	subtaskID := parts[6]

	if err := h.tasks.ToggleSubtask(subtaskID, taskID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *APIHandler) SubtaskDelete(w http.ResponseWriter, r *http.Request) {
	// Path: /api/v1/tasks/{task_id}/subtasks/{subtask_id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	subtaskID := parts[6]

	if err := h.tasks.DeleteSubtask(subtaskID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Checklists
func (h *APIHandler) ChecklistCreate(w http.ResponseWriter, r *http.Request) {
	taskID := h.extractID(r.URL.Path, "/api/v1/tasks/")
	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		title = "Checklist"
	}

	c, err := h.kanban.CreateChecklist(taskID, title, 1000.0)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, c)
}

func (h *APIHandler) ChecklistDelete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	taskID := parts[4]
	checklistID := parts[6]

	if err := h.kanban.DeleteChecklist(taskID, checklistID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Checklist Items
func (h *APIHandler) ChecklistItemCreate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 7 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	checklistID := parts[6]

	var input struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	item, err := h.kanban.AddChecklistItem(checklistID, input.Title, 1000.0)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, item)
}

func (h *APIHandler) ChecklistItemUpdate(w http.ResponseWriter, r *http.Request) {
	// Path: /api/v1/tasks/{task_id}/checklists/{checklist_id}/items/{item_id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 9 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	itemID := parts[8]

	var input struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	item := &models.ChecklistItem{
		ID:        itemID,
		Title:     input.Title,
		Completed: input.Completed,
	}

	if err := h.kanban.UpdateChecklistItem(item); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, item)
}

func (h *APIHandler) ChecklistItemDelete(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 9 {
		h.writeError(w, http.StatusBadRequest, "Invalid URL path")
		return
	}
	itemID := parts[8]

	if err := h.kanban.DeleteChecklistItem(itemID); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Comments & Attachments & Labels
func (h *APIHandler) CommentCreate(w http.ResponseWriter, r *http.Request) {
	taskID := h.extractID(r.URL.Path, "/api/v1/tasks/")
	var input struct {
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	c, err := h.kanban.CreateComment(taskID, input.Content)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, c)
}

func (h *APIHandler) AttachmentCreate(w http.ResponseWriter, r *http.Request) {
	taskID := h.extractID(r.URL.Path, "/api/v1/tasks/")
	var input struct {
		Filename string `json:"filename"`
		URL      string `json:"url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	a, err := h.kanban.CreateAttachment(taskID, input.Filename, input.URL, 100)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, a)
}

func (h *APIHandler) LabelToggle(w http.ResponseWriter, r *http.Request) {
	taskID := h.extractID(r.URL.Path, "/api/v1/tasks/")
	var input struct {
		LabelID string `json:"label_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	labels, _ := h.kanban.GetTaskLabels(taskID)
	exists := false
	for _, l := range labels {
		if l.ID == input.LabelID {
			exists = true
			break
		}
	}

	var err error
	if exists {
		err = h.kanban.RemoveLabelFromTask(taskID, input.LabelID)
	} else {
		err = h.kanban.AddLabelToTask(taskID, input.LabelID)
	}

	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *APIHandler) LabelList(w http.ResponseWriter, r *http.Request) {
	labels, err := h.kanban.ListLabels()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, labels)
}

// Notes
func (h *APIHandler) NoteList(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	notes, err := h.notes.Search(q)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, notes)
}

func (h *APIHandler) NoteCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		h.writeError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	n, err := h.notes.Create(title, input.Content)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, n)
}

func (h *APIHandler) NoteUpdate(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/notes/")
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	title := strings.TrimSpace(input.Title)
	if title == "" {
		h.writeError(w, http.StatusBadRequest, "Title cannot be empty")
		return
	}

	n, err := h.notes.Update(id, title, input.Content)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, n)
}

func (h *APIHandler) NoteDelete(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/notes/")
	if err := h.notes.Delete(id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

// Timeline
func (h *APIHandler) TimelineList(w http.ResponseWriter, r *http.Request) {
	acts, err := h.activities.List(100)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, acts)
}

// Timer
func (h *APIHandler) TimerStart(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TaskID      string `json:"task_id"`
		SubtaskID   string `json:"subtask_id"`
		SessionType string `json:"session_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	var pTaskID, pSubtaskID *string
	if input.TaskID != "" {
		pTaskID = &input.TaskID
	}
	if input.SubtaskID != "" {
		pSubtaskID = &input.SubtaskID
	}
	sessionType := input.SessionType
	if sessionType == "" {
		sessionType = "pomodoro"
	}

	session, err := h.timelogs.StartSession(pTaskID, pSubtaskID, sessionType)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, session)
}

func (h *APIHandler) TimerStop(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/timer/")
	id = strings.TrimSuffix(id, "/stop")

	var input struct {
		DurationSeconds int `json:"duration_seconds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	if err := h.timelogs.StopSession(id, input.DurationSeconds); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *APIHandler) TimerActive(w http.ResponseWriter, r *http.Request) {
	session, err := h.timelogs.GetActiveSession()
	if err != nil {
		h.writeJSON(w, http.StatusOK, nil)
		return
	}
	h.writeJSON(w, http.StatusOK, session)
}

// Reports
func (h *APIHandler) ReportList(w http.ResponseWriter, r *http.Request) {
	logs, err := h.timelogs.ListAllWithDetails()
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, logs)
}

// Reminders
func (h *APIHandler) ReminderCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TaskID   string `json:"task_id"`
		Note     string `json:"note"`
		RemindAt string `json:"remind_at"` // RFC3339 or "YYYY-MM-DD HH:MM"
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON input")
		return
	}

	remindAt, err := time.Parse("2006-01-02 15:04", input.RemindAt)
	if err != nil {
		remindAt, err = time.Parse(time.RFC3339, input.RemindAt)
	}
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid remind_at format")
		return
	}

	rem, err := h.reminders.Create(input.TaskID, input.Note, remindAt)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusCreated, rem)
}

func (h *APIHandler) ReminderDone(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/reminders/")
	id = strings.TrimSuffix(id, "/done")

	if err := h.reminders.MarkDone(id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

func (h *APIHandler) ReminderDelete(w http.ResponseWriter, r *http.Request) {
	id := h.extractID(r.URL.Path, "/api/v1/reminders/")
	if err := h.reminders.Delete(id); err != nil {
		h.writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	h.writeJSON(w, http.StatusOK, map[string]bool{"success": true})
}

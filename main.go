package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"activity-monitor/internal/database"
	"activity-monitor/internal/handlers"
	"activity-monitor/internal/repositories"
	"activity-monitor/internal/services"
)

//go:embed static
var staticFiles embed.FS

func main() {
	// Init database
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "activity-monitor.db"
	}

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	// Repositories
	taskRepo := repositories.NewTaskRepository(db)
	subtaskRepo := repositories.NewSubtaskRepository(db)
	reminderRepo := repositories.NewReminderRepository(db)
	noteRepo := repositories.NewNoteRepository(db)
	activityRepo := repositories.NewActivityRepository(db)
	timelogRepo := repositories.NewTimeLogRepository(db)
	kanbanRepo := repositories.NewKanbanRepository(db)

	// Services
	taskSvc := services.NewTaskService(taskRepo, subtaskRepo, reminderRepo, activityRepo, timelogRepo, kanbanRepo)
	noteSvc := services.NewNoteService(noteRepo, activityRepo)
	reminderSvc := services.NewReminderService(reminderRepo, taskRepo, activityRepo)
	activitySvc := services.NewActivityService(activityRepo)
	timelogSvc := services.NewTimeLogService(timelogRepo)
	kanbanSvc := services.NewKanbanService(kanbanRepo, taskRepo, activityRepo)

	// Handlers
	dashboardH := handlers.NewDashboardHandler(taskSvc, reminderSvc, activitySvc, timelogSvc)
	taskH := handlers.NewTaskHandler(taskSvc, reminderSvc)
	noteH := handlers.NewNoteHandler(noteSvc)
	timelineH := handlers.NewTimelineHandler(activitySvc)
	reminderH := handlers.NewReminderHandler(reminderSvc)
	timerH := handlers.NewTimerHandler(timelogSvc)
	reportH := handlers.NewReportHandler(timelogSvc)
	kanbanH := handlers.NewKanbanHandler(kanbanSvc, taskSvc)

	mux := http.NewServeMux()

	// Static files - embedded from ./static
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Ignore favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// Dashboard
	mux.HandleFunc("/", dashboardH.Show)
	mux.HandleFunc("/components/sidebar-stats", dashboardH.SidebarStats)

	// Tasks
	mux.HandleFunc("/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			taskH.List(w, r)
		case http.MethodPost:
			taskH.Create(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/tasks/new", taskH.NewForm)

	mux.HandleFunc("/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.HasSuffix(path, "/edit") && r.Method == http.MethodGet:
			taskH.EditForm(w, r)
		case strings.HasSuffix(path, "/status") && r.Method == http.MethodPut:
			taskH.UpdateStatus(w, r)
		case strings.HasSuffix(path, "/move") && r.Method == http.MethodPut:
			kanbanH.MoveTask(w, r)
		case strings.HasSuffix(path, "/comments") && r.Method == http.MethodPost:
			kanbanH.AddComment(w, r)
		case strings.HasSuffix(path, "/checklists") && r.Method == http.MethodPost:
			kanbanH.AddChecklist(w, r)
		case strings.HasSuffix(path, "/attachments") && r.Method == http.MethodPost:
			kanbanH.AddAttachment(w, r)
		case strings.HasSuffix(path, "/labels") && r.Method == http.MethodPost:
			kanbanH.ToggleLabel(w, r)
		case strings.HasSuffix(path, "/pomodoro") && r.Method == http.MethodPost:
			kanbanH.UpdatePomodoro(w, r)
		case strings.HasSuffix(path, "/subtasks") && r.Method == http.MethodPost:
			taskH.AddSubtask(w, r)
		case strings.Contains(path, "/subtasks/") && strings.HasSuffix(path, "/toggle") && r.Method == http.MethodPut:
			taskH.ToggleSubtask(w, r)
		case strings.Contains(path, "/subtasks/") && r.Method == http.MethodDelete:
			taskH.DeleteSubtask(w, r)
		case strings.HasSuffix(path, "/reminder/new") && r.Method == http.MethodGet:
			taskH.ReminderForm(w, r)
		case strings.HasSuffix(path, "/reminders") && r.Method == http.MethodPost:
			taskH.CreateReminder(w, r)
		case r.Method == http.MethodPut:
			taskH.Update(w, r)
		case r.Method == http.MethodDelete:
			taskH.Delete(w, r)
		case r.Method == http.MethodGet:
			taskH.Detail(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// Notes
	mux.HandleFunc("/notes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			noteH.List(w, r)
		case http.MethodPost:
			noteH.Create(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	mux.HandleFunc("/notes/new", noteH.NewForm)
	mux.HandleFunc("/notes/search", noteH.Search)

	mux.HandleFunc("/notes/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/edit") && r.Method == http.MethodGet:
			noteH.EditForm(w, r)
		case r.Method == http.MethodPut:
			noteH.Update(w, r)
		case r.Method == http.MethodDelete:
			noteH.Delete(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	// Timeline
	mux.HandleFunc("/timeline", timelineH.Show)

	// Reports
	mux.HandleFunc("/reports", reportH.Show)

	// Timer
	mux.HandleFunc("/timer/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			timerH.Start(w, r)
		}
	})
	mux.HandleFunc("/timer/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/stop") && r.Method == http.MethodPost {
			timerH.Stop(w, r)
		}
	})

	// Reminders
	mux.HandleFunc("/reminders/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/done") && r.Method == http.MethodPut:
			reminderH.MarkDone(w, r)
		case r.Method == http.MethodDelete:
			reminderH.Delete(w, r)
		default:
			http.NotFound(w, r)
		}
	})

	log.Println("Activity Monitor berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

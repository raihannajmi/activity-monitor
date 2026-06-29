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
	apiH := handlers.NewAPIHandler(taskSvc, noteSvc, reminderSvc, activitySvc, timelogSvc, kanbanSvc)

	mux := http.NewServeMux()

	// Static files - embedded from ./static
	staticFS, _ := fs.Sub(staticFiles, "static")
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

	// Ignore favicon
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// API Routing
	mux.HandleFunc("/api/v1/dashboard", apiH.DashboardStats)
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			apiH.TaskCreate(w, r)
		} else {
			apiH.TaskList(w, r)
		}
	})
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		switch {
		case strings.HasSuffix(path, "/move"):
			apiH.TaskMove(w, r)
		case strings.HasSuffix(path, "/subtasks"):
			apiH.SubtaskCreate(w, r)
		case strings.Contains(path, "/subtasks/") && strings.HasSuffix(path, "/toggle"):
			apiH.SubtaskToggle(w, r)
		case strings.Contains(path, "/subtasks/"):
			apiH.SubtaskDelete(w, r)
		case strings.HasSuffix(path, "/checklists"):
			apiH.ChecklistCreate(w, r)
		case strings.Contains(path, "/checklists/") && strings.Contains(path, "/items/"):
			if r.Method == http.MethodDelete {
				apiH.ChecklistItemDelete(w, r)
			} else {
				apiH.ChecklistItemUpdate(w, r)
			}
		case strings.Contains(path, "/checklists/") && strings.HasSuffix(path, "/items"):
			apiH.ChecklistItemCreate(w, r)
		case strings.Contains(path, "/checklists/"):
			apiH.ChecklistDelete(w, r)
		case strings.HasSuffix(path, "/comments"):
			apiH.CommentCreate(w, r)
		case strings.HasSuffix(path, "/attachments"):
			apiH.AttachmentCreate(w, r)
		case strings.HasSuffix(path, "/labels"):
			apiH.LabelToggle(w, r)
		default:
			if r.Method == http.MethodDelete {
				apiH.TaskDelete(w, r)
			} else if r.Method == http.MethodPut {
				apiH.TaskUpdate(w, r)
			} else {
				apiH.TaskGet(w, r)
			}
		}
	})
	mux.HandleFunc("/api/v1/labels", apiH.LabelList)
	mux.HandleFunc("/api/v1/notes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			apiH.NoteCreate(w, r)
		} else {
			apiH.NoteList(w, r)
		}
	})
	mux.HandleFunc("/api/v1/notes/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			apiH.NoteDelete(w, r)
		} else {
			apiH.NoteUpdate(w, r)
		}
	})
	mux.HandleFunc("/api/v1/timeline", apiH.TimelineList)
	mux.HandleFunc("/api/v1/timer/start", apiH.TimerStart)
	mux.HandleFunc("/api/v1/timer/active", apiH.TimerActive)
	mux.HandleFunc("/api/v1/timer/", apiH.TimerStop)
	mux.HandleFunc("/api/v1/reports", apiH.ReportList)
	mux.HandleFunc("/api/v1/reminders", apiH.ReminderCreate)
	mux.HandleFunc("/api/v1/reminders/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/done") {
			apiH.ReminderDone(w, r)
		} else {
			apiH.ReminderDelete(w, r)
		}
	})

	// Dashboard / templ legacy views (kept for compatibility during build)
	mux.HandleFunc("/legacy-dashboard", dashboardH.Show)
	mux.HandleFunc("/components/sidebar-stats", dashboardH.SidebarStats)


	// Root handler (serves static/index.html for SPA routes)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// If requesting api, return 404
		if strings.HasPrefix(r.URL.Path, "/api/") {
			http.NotFound(w, r)
			return
		}
		// If requesting static assets, or if path has file extension, it should not serve index.html
		if strings.HasPrefix(r.URL.Path, "/static/") || strings.Contains(r.URL.Path, ".") {
			http.NotFound(w, r)
			return
		}

		// Read index.html from staticFiles
		content, err := staticFiles.ReadFile("static/index.html")
		if err != nil {
			// Fallback if index.html doesn't exist yet (during initial dev build)
			w.Header().Set("Content-Type", "text/html")
			w.Write([]byte("<h1>Activity Monitor Vue SPA loading... (Please build the frontend)</h1>"))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(content)
	})

	log.Println("Activity Monitor berjalan di http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server: %v", err)
	}
}

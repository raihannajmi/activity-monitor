package handlers

import (
	"net/http"
	"strings"

	"activity-monitor/internal/services"
	"activity-monitor/templates/components"
	"activity-monitor/templates/pages"
)

type NoteHandler struct {
	notes *services.NoteService
}

func NewNoteHandler(notes *services.NoteService) *NoteHandler {
	return &NoteHandler{notes}
}

func (h *NoteHandler) List(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	noteList, err := h.notes.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Notes(noteList, query).Render(r.Context(), w)
}

func (h *NoteHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	noteList, err := h.notes.Search(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.NotesList(noteList).Render(r.Context(), w)
}

func (h *NoteHandler) NewForm(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.NoteFormModal(nil).Render(r.Context(), w)
}

func (h *NoteHandler) EditForm(w http.ResponseWriter, r *http.Request) {
	id := extractNoteID(r.URL.Path)
	note, err := h.notes.GetByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	components.NoteFormModal(note).Render(r.Context(), w)
}

func (h *NoteHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "judul tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if _, err := h.notes.Create(title, r.FormValue("content")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteList, _ := h.notes.List()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.NotesList(noteList).Render(r.Context(), w)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractNoteID(r.URL.Path)
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "judul tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if _, err := h.notes.Update(id, title, r.FormValue("content")); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	noteList, _ := h.notes.List()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.NotesList(noteList).Render(r.Context(), w)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractNoteID(r.URL.Path)
	if err := h.notes.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func extractNoteID(path string) string {
	path = strings.TrimPrefix(path, "/notes/")
	if idx := strings.Index(path, "/"); idx != -1 {
		path = path[:idx]
	}
	return path
}

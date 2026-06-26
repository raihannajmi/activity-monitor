package handlers

import (
	"net/http"
	"net/url"
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
		h.sendError(w, r, "Gagal mengambil daftar catatan", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	pages.Notes(noteList, query).Render(r.Context(), w)
}

func (h *NoteHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	noteList, err := h.notes.Search(query)
	if err != nil {
		h.sendError(w, r, "Gagal melakukan pencarian", http.StatusInternalServerError)
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
		h.sendError(w, r, "Data tidak valid", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		h.sendError(w, r, "Judul catatan tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if _, err := h.notes.Create(title, r.FormValue("content")); err != nil {
		h.sendError(w, r, "Gagal membuat catatan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", r.Header.Get("HX-Current-URL"))
	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := extractNoteID(r.URL.Path)
	if err := r.ParseForm(); err != nil {
		h.sendError(w, r, "Data tidak valid", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		h.sendError(w, r, "Judul catatan tidak boleh kosong", http.StatusBadRequest)
		return
	}

	if _, err := h.notes.Update(id, title, r.FormValue("content")); err != nil {
		h.sendError(w, r, "Gagal menyimpan catatan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Location", r.Header.Get("HX-Current-URL"))
	w.WriteHeader(http.StatusOK)
}

func (h *NoteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := extractNoteID(r.URL.Path)
	if err := h.notes.Delete(id); err != nil {
		h.sendError(w, r, "Gagal menghapus catatan", http.StatusInternalServerError)
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

func (h *NoteHandler) extractQuery(r *http.Request) string {
	currentURL := r.Header.Get("HX-Current-URL")
	if currentURL == "" {
		return ""
	}
	if u, err := url.Parse(currentURL); err == nil {
		return u.Query().Get("q")
	}
	return ""
}

func (h *NoteHandler) sendError(w http.ResponseWriter, r *http.Request, msg string, statusCode int) {
	w.Header().Set("HX-Retarget", "#toast-container")
	w.Header().Set("HX-Reswap", "beforeend")
	w.WriteHeader(http.StatusOK)
	components.ToastError(msg).Render(r.Context(), w)
}

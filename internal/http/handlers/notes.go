package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"example.com/goprac11-borisovda/internal/core"
	"example.com/goprac11-borisovda/internal/repo"
	"github.com/go-chi/chi/v5"
)

type Handler struct {
	Repo *repo.NoteRepoMem
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var n core.Note
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	id, _ := h.Repo.Create(n)
	n.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(n)
}

func (h *Handler) GetAllNotes(w http.ResponseWriter, r *http.Request) {
	notes, _ := h.Repo.GetAll()
	writeJSON(w, http.StatusOK, notes)
}

func (h *Handler) GetNote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	n, _ := h.Repo.Get(id)
	if n == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	var upd core.Note
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	n, _ := h.Repo.Update(id, upd)
	if n == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)

	h.Repo.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}

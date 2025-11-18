package handler

import (
	"bookapi/internal/model"
	"bookapi/internal/store"
	"encoding/json"
	"net/http"
	"strings"
)

type BookHandler struct {
	store *store.ShardedBookStore
}

func NewBookHandler(store *store.ShardedBookStore) *BookHandler {
	return &BookHandler{store: store}
}

func respondJson(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJson(w, status, map[string]string{"error": message})
}

func extractID(path, prefix string) string {
	return strings.TrimPrefix(path, prefix)
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := h.store.GetAll()
	respondJson(w, http.StatusOK, map[string]interface{}{
		"books": books,
		"count": len(books),
	})
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/books/")
	book, ok := h.store.GetById(id)
	if !ok {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}
	respondJson(w, http.StatusOK, book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req model.CreateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if req.Title == "" || req.Author == "" {
		respondError(w, http.StatusBadRequest, "Title and Author are required")
	}
	book := h.store.Create(req)
	respondJson(w, http.StatusCreated, book)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/books/")
	var req model.UpdateBookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	book, ok := h.store.Update(id, req)
	if !ok {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}
	respondJson(w, http.StatusOK, book)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	id := extractID(r.URL.Path, "/books/")
	if !h.store.Delete(id) {
		respondError(w, http.StatusNotFound, "Book not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

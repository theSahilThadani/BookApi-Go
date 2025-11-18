package handler

import (
	"bookapi/internal/model"
	"bookapi/internal/store"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// helper to create handler with seeded data
func newTestHandler() *BookHandler {
	s := store.NewShardedBookStore(22)
	s.Create(model.CreateBookRequest{Title: "T1", Author: "A1", ISBN: "X1", PublishYear: 2020})
	return NewBookHandler(s)
}

func TestGetBooksHandler(t *testing.T) {
	h := newTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	w := httptest.NewRecorder()

	h.GetBooks(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body error: %v", err)
	}
	if _, ok := body["books"]; !ok {
		t.Fatalf("expected books key in response")
	}
}

func TestCreateBookHandler(t *testing.T) {
	h := newTestHandler()

	reqBody := model.CreateBookRequest{
		Title: "New Book", Author: "Author", ISBN: "ISBN-1", PublishYear: 2024,
	}
	b, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateBook(w, req)

	if w.Result().StatusCode != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Result().StatusCode)
	}

	var created model.Book
	if err := json.NewDecoder(w.Result().Body).Decode(&created); err != nil {
		t.Fatalf("decode created book error: %v", err)
	}
	if created.Title != reqBody.Title {
		t.Fatalf("title mismatch: want %s got %s", reqBody.Title, created.Title)
	}
}

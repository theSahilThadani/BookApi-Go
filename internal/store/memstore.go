package store

import (
	"bookapi/internal/model"
	"fmt"
	"sync"
	"time"
)

type BookStore struct {
	mu     sync.RWMutex
	books  map[string]*model.Book
	nextID int
}

func NewBookStore() *BookStore {
	return &BookStore{
		books:  make(map[string]*model.Book),
		nextID: 1,
	}
}

func (s *BookStore) Create(req model.CreateBookRequest) *model.Book {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("BOOK-%d", s.nextID)
	s.nextID++

	book := &model.Book{
		ID:          id,
		Title:       req.Title,
		Author:      req.Author,
		ISBN:        req.ISBN,
		PublishYear: req.PublishYear,
		CreatedAt:   time.Now(),
	}
	s.books[id] = book
	return book
}
func (s *BookStore) GetAll() []*model.Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]*model.Book, 0, len(s.books))
	for _, b := range s.books {
		out = append(out, b)
	}
	return out
}

func (s *BookStore) GetById(id string) (*model.Book, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.books[id]
	return b, ok
}

func (s *BookStore) Update(id string, req model.UpdateBookRequest) (*model.Book, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	book, ok := s.books[id]
	if !ok {
		return nil, false
	}
	if req.Title != nil {
		book.Title = *req.Title
	}
	if req.Author != nil {
		book.Author = *req.Author
	}
	if req.ISBN != nil {
		book.ISBN = *req.ISBN
	}
	if req.PublishYear != nil {
		book.PublishYear = *req.PublishYear
	}

	return book, true
}
func (s *BookStore) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return false
	}
	delete(s.books, id)
	return true
}

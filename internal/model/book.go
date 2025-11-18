package model

import "time"

// Book represents a book in our library
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	ISBN        string    `json:"isbn"`
	PublishYear int       `json:"publish_year"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateBookRequest - used when creating a new book (required fields)
type CreateBookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	ISBN        string `json:"isbn"`
	PublishYear int    `json:"publish_year"`
}

// UpdateBookRequest - used when updating a book; pointers -> optional fields
type UpdateBookRequest struct {
	Title       *string `json:"title,omitempty"`
	Author      *string `json:"author,omitempty"`
	ISBN        *string `json:"isbn,omitempty"`
	PublishYear *int    `json:"publish_year,omitempty"`
}

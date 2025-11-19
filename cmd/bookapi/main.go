package main

import (
	"bookapi/internal/handler"
	"bookapi/internal/middleware"
	"bookapi/internal/model"
	"bookapi/internal/store"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {

	go func() {
		log.Println("pprof at :6060")
		http.ListenAndServe(":6060", nil)
	}()

	// init store and seed sample data
	s := store.NewShardedBookStore(22)
	s.Create(model.CreateBookRequest(storeCreate("The Go Programming Language", "Alan Donovan & Brian Kernighan", "978-0134190440", 2015)))
	s.Create(model.CreateBookRequest(storeCreate("Clean Code", "Robert C. Martin", "978-0132350884", 2008)))

	h := handler.NewBookHandler(s)

	mux := http.NewServeMux()

	// Route: /books  list and create
	mux.Handle("/books", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetBooks(w, r)
		case http.MethodPost:
			h.CreateBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	// Route: /books/  get/update/delete by id
	mux.Handle("/books/", middleware.Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.GetBook(w, r)
		case http.MethodPut:
			h.UpdateBook(w, r)
		case http.MethodDelete:
			h.DeleteBook(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	addr := ":8080"
	log.Printf("Server starting at http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func storeCreate(title, author, isbn string, year int) storeCreateReq {
	return storeCreateReq{title, author, isbn, year}
}

type storeCreateReq struct {
	Title       string
	Author      string
	ISBN        string
	PublishYear int
}

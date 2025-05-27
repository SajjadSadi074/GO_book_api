package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"bookapi/models"
	"github.com/go-chi/chi/v5"
)

type BookHandler struct {
	library *models.Library
}

// OOP: Dependency Injection
func NewBookHandler(lib *models.Library) *BookHandler {
	return &BookHandler{library: lib}
}

func (h *BookHandler) Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Book API!"))
}

func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.library.GetBooks())
}

func (h *BookHandler) GetBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	book, err := h.library.GetBookByISBN(isbn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

func (h *BookHandler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if newBook.Title == "" || newBook.Author == "" || newBook.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	if err := h.library.AddBook(newBook); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func (h *BookHandler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	var updatedBook models.Book
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &updatedBook)

	_, err := h.library.GetBookByISBN(isbn)
	if err != nil {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	if isbn != updatedBook.ISBN {
		if _, err := h.library.GetBookByISBN(updatedBook.ISBN); err == nil {
			http.Error(w, "Book with this ISBN already exists", http.StatusConflict)
			return
		}
	}
	_ = h.library.DeleteBookByISBN(isbn)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	h.CreateBook(w, r)
}

func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	if err := h.library.DeleteBookByISBN(isbn); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

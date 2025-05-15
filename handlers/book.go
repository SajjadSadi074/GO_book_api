package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"bookapi/models"
	"github.com/go-chi/chi/v5"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Book API!"))
}

func GetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Books)
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	for _, book := range models.Books {
		if book.ISBN == isbn {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var newBook models.Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if newBook.Title == "" || newBook.Author == "" || newBook.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	for _, b := range models.Books {
		if b.ISBN == newBook.ISBN {
			http.Error(w, "Book with this ISBN already exists", http.StatusConflict)
			return
		}
	}
	models.Books = append(models.Books, newBook)
	models.AddAuthor(newBook.Author)
	models.AddAuthorBook(newBook.ISBN, newBook.Author)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	var updatedBook models.Book
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &updatedBook)

	var found bool
	for _, b := range models.Books {
		if b.ISBN == isbn {
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	if updatedBook.Title == "" || updatedBook.Author == "" || updatedBook.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}
	if isbn != updatedBook.ISBN {
		for _, b := range models.Books {
			if b.ISBN == updatedBook.ISBN {
				http.Error(w, "Book with this ISBN already exists", http.StatusConflict)
				return
			}
		}
	}
	DeleteBook(w, r)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	CreateBook(w, r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBook)
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	var index int
	var found bool
	for i, b := range models.Books {
		if b.ISBN == isbn {
			index = i
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "Book not found", http.StatusNotFound)
		return
	}

	author := strings.ToLower(models.Books[index].Author)
	for i, v := range models.AuthorBooks[author] {
		if v == isbn {
			models.AuthorBooks[author] = append(models.AuthorBooks[author][:i], models.AuthorBooks[author][i+1:]...)
			break
		}
	}
	models.Books = append(models.Books[:index], models.Books[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

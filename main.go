package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"io"
	"log"
	"net/http"
	"strings"
)

var tokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

var books = []Book{
	{Title: "Go Programming", Author: "Alice", ISBN: "123-ABC"},
	{Title: "Microservices in Go", Author: "Bob", ISBN: "456-DEF"},
}

const (
	username = "admin"
	password = "secret"
)

var authors = []string{
	"Alice",
	"Bob",
}

var authorBooks = map[string][]string{
	"alice": {"123-ABC"},
	"bob":   {"456-DEF"},
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// If the auther is already there, then it doesn't add
func add_author(str string) {
	if contains(authors, str) {
		return
	}
	authors = append(authors, str)
}

// If the book is already there than it doesn't add,
// it useses isbn to determine uniqness
func add_author_book(isbn string, author string) {
	author = strings.ToLower(author)
	if contains(authorBooks[author], isbn) {
		return
	}
	authorBooks[author] = append(authorBooks[author], isbn)
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Book API!"))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	for _, book := range books {
		if book.ISBN == isbn {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(book)
			return
		}
	}
	http.Error(w, "Book not found", http.StatusNotFound)
}

// print all the authors available
func getAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authors)
}

// prints every book infomation of a perticular writer
func getAuthorBook(w http.ResponseWriter, r *http.Request) {
	authorName := chi.URLParam(r, "author")
	fmt.Println(authorName)
	authorName = strings.ToLower(authorName)
	EveryBookofAuther := authorBooks[authorName]
	PrintBooks := []Book{} // PrintBooks is used for temporarily storing the data that needs to be outputed

	//This nested loop is used for matching every isbn for
	//the author with every book's isbn,

	//The Compexity can be improved by using a map, where the books are stored directly for every author,
	//or a map where the book information is stored for every isbn
	for _, book := range books {
		for _, isbn := range EveryBookofAuther {
			if book.ISBN == isbn {
				PrintBooks = append(PrintBooks, book)
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(PrintBooks)
	return
}

func createBook(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Creating new book=====")
	var newBook Book
	if err := json.NewDecoder(r.Body).Decode(&newBook); err != nil {
		fmt.Println(err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	fmt.Println(newBook)

	if newBook.Title == "" || newBook.Author == "" || newBook.ISBN == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	for _, b := range books {
		if b.ISBN == newBook.ISBN {
			http.Error(w, "Book with this ISBN already exists", http.StatusConflict)
			return
		}
	}

	books = append(books, newBook)
	add_author(newBook.Author)
	add_author_book(newBook.ISBN, newBook.Author)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newBook)
}

// Update works using delete and create function. Before deleting anything, first we check is it
// creatable or not, to ensure that data is not lost.
func updateBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	var updatedBook Book
	bodyBytes, _ := io.ReadAll(r.Body)
	json.Unmarshal(bodyBytes, &updatedBook)

	var found bool
	for _, b := range books {
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
		for _, b := range books {
			if b.ISBN == updatedBook.ISBN {
				http.Error(w, "Book with this ISBN already exists", http.StatusConflict)
				return
			}
		}
	}

	deleteBook(w, r)
	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	createBook(w, r)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedBook)
}

func deleteBook(w http.ResponseWriter, r *http.Request) {
	isbn := chi.URLParam(r, "isbn")
	var index int
	var found bool
	for i, b := range books {
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

	//deleted from author books
	var tmp = books[index].Author
	tmp = strings.ToLower(tmp)
	var index_author_book int
	for i, value := range authorBooks[tmp] {
		if value == books[index].ISBN {
			index_author_book = i
			break
		}
	}
	fmt.Println(index_author_book)
	authorBooks[tmp] = append(authorBooks[tmp][:index_author_book], authorBooks[tmp][index_author_book+1:]...)

	//deleted from books
	books = append(books[:index], books[index+1:]...)

	w.WriteHeader(http.StatusNoContent)
}

func basicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || user != username || pass != password {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func login(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")

	// In real-world use, validate these credentials from a database
	if username != "admin" || password != "pass123" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Create JWT token
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user": username})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func main() {
	r := chi.NewRouter()

	// Global middleware – these must come first
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// JWT setup
	tokenAuth := jwtauth.New("HS256", []byte("your-secret-key"), nil)

	// Public route for login
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		_, token, _ := tokenAuth.Encode(map[string]interface{}{"user_id": 123})
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// JWT middleware applies only inside this group
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/books", createBook)
		r.Put("/books/{isbn}", updateBook)
		r.Delete("/books/{isbn}", deleteBook)
	})

	// Unprotected routes
	r.Get("/books", getBooks)
	r.Get("/", home)
	r.Get("/books/{isbn}", getBook)
	r.Get("/authors", getAuthors)
	r.Get("/authors/{author}", getAuthorBook)
	r.Get("/login", login)

	log.Println("Starting server on http://127.0.0.1:8080...")
	if err := http.ListenAndServe("127.0.0.1:8080", r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

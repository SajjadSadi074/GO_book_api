package models

import (
	"errors"
	"strings"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

type Library struct {
	books       []Book
	authors     []string
	authorBooks map[string][]string
}

func NewLibrary() *Library {
	return &Library{
		books: []Book{
			{Title: "Go Programming", Author: "Alice", ISBN: "123-ABC"},
			{Title: "Microservices in Go", Author: "Bob", ISBN: "456-DEF"},
		},
		authors: []string{"Alice", "Bob"},
		authorBooks: map[string][]string{
			"alice": {"123-ABC"},
			"bob":   {"456-DEF"},
		},
	}
}

// OOP: Encapsulation (data and methods inside Library)
func (l *Library) GetBooks() []Book {
	return l.books
}

func (l *Library) GetBookByISBN(isbn string) (*Book, error) {
	for _, b := range l.books {
		if b.ISBN == isbn {
			return &b, nil
		}
	}
	return nil, errors.New("book not found")
}

func (l *Library) AddBook(book Book) error {
	for _, b := range l.books {
		if b.ISBN == book.ISBN {
			return errors.New("book with this ISBN already exists")
		}
	}
	l.books = append(l.books, book)
	l.AddAuthor(book.Author)
	l.AddAuthorBook(book.ISBN, book.Author)
	return nil
}

func (l *Library) DeleteBookByISBN(isbn string) error {
	for i, b := range l.books {
		if b.ISBN == isbn {
			author := strings.ToLower(b.Author)
			// Remove from authorBooks map
			if books, ok := l.authorBooks[author]; ok {
				for j, bookISBN := range books {
					if bookISBN == isbn {
						l.authorBooks[author] = append(books[:j], books[j+1:]...)
						break
					}
				}
			}
			// Remove from books slice
			l.books = append(l.books[:i], l.books[i+1:]...)
			return nil
		}
	}
	return errors.New("book not found")
}

func (l *Library) AddAuthor(author string) {
	if l.contains(l.authors, author) {
		return
	}
	l.authors = append(l.authors, author)
}

func (l *Library) AddAuthorBook(isbn string, author string) {
	author = strings.ToLower(author)
	if l.contains(l.authorBooks[author], isbn) {
		return
	}
	l.authorBooks[author] = append(l.authorBooks[author], isbn)
}

func (l *Library) contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func (l *Library) GetAuthors() []string {
	return l.authors
}

func (l *Library) GetBooksByAuthor(author string) []Book {
	author = strings.ToLower(author)
	isbns := l.authorBooks[author]
	var books []Book

	// Optional optimization: build an ISBN map for O(1) lookup
	isbnMap := make(map[string]Book)
	for _, b := range l.books {
		isbnMap[b.ISBN] = b
	}

	for _, isbn := range isbns {
		if book, ok := isbnMap[isbn]; ok {
			books = append(books, book)
		}
	}

	return books
}

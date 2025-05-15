package models

import "strings"

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	ISBN   string `json:"isbn"`
}

var Books = []Book{
	{Title: "Go Programming", Author: "Alice", ISBN: "123-ABC"},
	{Title: "Microservices in Go", Author: "Bob", ISBN: "456-DEF"},
}

var Authors = []string{
	"Alice",
	"Bob",
}

const (
	username = "admin"
	password = "secret"
)

var AuthorBooks = map[string][]string{
	"alice": {"123-ABC"},
	"bob":   {"456-DEF"},
}

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func AddAuthor(str string) {
	if Contains(Authors, str) {
		return
	}
	Authors = append(Authors, str)
}

func AddAuthorBook(isbn string, author string) {
	author = strings.ToLower(author)
	if Contains(AuthorBooks[author], isbn) {
		return
	}
	AuthorBooks[author] = append(AuthorBooks[author], isbn)
}

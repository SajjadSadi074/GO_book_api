package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"bookapi/models"
	"github.com/go-chi/chi/v5"
)

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(models.Authors)
}

// prints every book infomation of a perticular writer
func GetAuthorBook(w http.ResponseWriter, r *http.Request) {
	authorName := chi.URLParam(r, "author")
	fmt.Println(authorName)
	authorName = strings.ToLower(authorName)
	EveryBookofAuther := models.AuthorBooks[authorName]
	PrintBooks := []models.Book{} // PrintBooks is used for temporarily storing the data that needs to be outputed

	//This nested loop is used for matching every isbn for
	//the author with every book's isbn,

	//The Compexity can be improved by using a map, where the books are stored directly for every author,
	//or a map where the book information is stored for every isbn
	for _, book := range models.Books {
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

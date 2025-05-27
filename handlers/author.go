package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *BookHandler) GetAuthors(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(h.library.GetAuthors())
}

func (h *BookHandler) GetAuthorBook(w http.ResponseWriter, r *http.Request) {
	author := chi.URLParam(r, "author")
	books := h.library.GetBooksByAuthor(author)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

package routes

import (
	"net/http"

	"bookapi/handlers"
	"bookapi/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)

func Routes() http.Handler {
	// Create the library and inject it into the handler
	library := models.NewLibrary()
	handler := handlers.NewBookHandler(library)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", handlers.Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/books", handler.CreateBook)
		r.Put("/books/{isbn}", handler.UpdateBook)
		r.Delete("/books/{isbn}", handler.DeleteBook)
	})

	r.Get("/", handler.Home)
	r.Get("/books", handler.GetBooks)
	r.Get("/books/{isbn}", handler.GetBook)
	r.Get("/authors", handler.GetAuthors)
	r.Get("/authors/{author}", handler.GetAuthorBook)

	return r
}

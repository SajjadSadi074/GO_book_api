package routes

import (
	"net/http"

	"bookapi/handlers" // use your module name here
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth = jwtauth.New("HS256", []byte("your-secret-key"), nil)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", handlers.Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/books", handlers.CreateBook)
		r.Put("/books/{isbn}", handlers.UpdateBook)
		r.Delete("/books/{isbn}", handlers.DeleteBook)
	})

	r.Get("/", handlers.Home)
	r.Get("/books", handlers.GetBooks)
	r.Get("/books/{isbn}", handlers.GetBook)
	r.Get("/authors", handlers.GetAuthors)
	r.Get("/authors/{author}", handlers.GetAuthorBook)

	return r
}

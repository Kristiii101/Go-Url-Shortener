package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// API routes
	r.Route("/api", func(api chi.Router) {
		api.Post("/shorten", shortenURL)
	})

	// Static files (dashboard)
	fileServer := http.FileServer(http.Dir("web"))
	r.Handle("/*", fileServer)

	return r
}

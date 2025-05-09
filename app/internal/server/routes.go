package server

import (
	"otus_social_network/internal/app/handlers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func ConfigureRouting() chi.Router {
	handlers := handlers.Init()
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/users/{id}", handlers.GetItem)

	return r
}

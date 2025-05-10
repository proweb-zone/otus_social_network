package server

import (
	"net/http"
	"otus_social_network/internal/app/handlers"
	"otus_social_network/internal/config"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

func StartServer(config *config.Config) {
	handlers := handlers.Init(config)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/login", handlers.Login)
	r.Put("/user/register", handlers.Register)
	r.Get("/user/get/{id}", handlers.GetUser)
	http.ListenAndServe(config.HTTPServer.ServerPort, r)
}

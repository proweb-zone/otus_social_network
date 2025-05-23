package server

import (
	"net/http"
	"otus_social_network/app/internal/app/handlers"
	"otus_social_network/app/internal/app/middleware"
	"otus_social_network/app/internal/config"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func StartServer(config *config.Config) {
	handlers := handlers.Init(config)

	r := chi.NewRouter()
	r.Post("/login", handlers.Login)
	r.Post("/user/register", handlers.Register)
	r.Get("/user/search/{query}", handlers.SearchUser)
	r.With(middleware.CheckAccess(config)).Get("/user/get/{id}", handlers.GetUser)
	http.ListenAndServe(":"+config.HTTPServer.ServerPort, r)
}

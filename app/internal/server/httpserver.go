package server

import (
	"net/http"
	"otus_social_network/app/internal/app/handlers"
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
	r.Get("/user/get/{id}", handlers.GetUser)

	r.Put("/friend/set/{user_id}", handlers.SetFriend)
	r.Put("/friend/delete/{user_id}", handlers.DelFriend)

	r.Post("/post/create", handlers.CreatePost)
	r.Put("/post/update", handlers.UpdatePost)
	r.Put("/post/delete/{id}", handlers.DelPost)
	r.Get("/post/get/{id}", handlers.GetPost)
	r.Get("/post/feed", handlers.FeedPost)
	http.ListenAndServe(":"+config.HTTPServer.ServerPort, r)
}

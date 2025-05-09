package server

import (
	"net/http"
	"time"

	"otus_social_network/internal/config"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func StartServer(router *chi.Mux, config *config.Config) *http.Server {

	server := &http.Server{
		Addr:         config.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  90 * time.Second,
	}

	return server
}

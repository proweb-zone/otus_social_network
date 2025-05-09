package app

import (
	"net/http"
	"otus_social_network/internal/server"
)

// Initialization
func InitApp() {
	router := server.ConfigureRouting()
	http.ListenAndServe(":3009", router)
}

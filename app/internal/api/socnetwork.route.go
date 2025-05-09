package api

import (
	"otus_social_network/internal/app/handlers"

	"github.com/go-chi/chi"
)

type SetupSocNetworkRoute struct {
	//logger *logger.Logger
}

func SetupSocNetworkRoutes() *SetupSocNetworkRoute {
	return &SetupSocNetworkRoute{}
}

func (route *SetupSocNetworkRoute) Init(r chi.Router) {
	socnetworkHandlers := handlers.Init()

	r.Get("/users/{id}", socnetworkHandlers.GetItem)
}

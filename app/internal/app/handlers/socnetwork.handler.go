package handlers

import (
	"fmt"
	"net/http"
)

type Handler struct {
	//service *service.SocNetworkService
}

func Init() *Handler {
	// socNetworkRepository := repository.InitPostgresRepository()
	// service := service.NewService(socNetworkRepository)

	return &Handler{}
}

func (c *Handler) GetItem(w http.ResponseWriter, r *http.Request) {
	//id := chi.URLParam(r, "id")

	fmt.Println("dfdf")
}

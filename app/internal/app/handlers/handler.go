package handlers

import (
	"net/http"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/app/service"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
)

type Handler struct {
	service *service.UserService
}

func Init(config *config.Config) *Handler {
	db := postgres.Connect(config.Db.StrConn)
	userRepository := repository.InitPostgresRepository(db)
	service := service.InitUserService(userRepository)

	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	h.service.Login()
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	h.service.Register()
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	h.service.GetUserById()
}

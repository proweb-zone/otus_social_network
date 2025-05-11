package handlers

import (
	"fmt"
	"io"
	"net/http"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/app/service"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
	"otus_social_network/app/internal/utils"
	"strconv"

	"github.com/go-chi/chi"
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
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestDto dto.UsersRequestDto
	if err := utils.DecodeJson(body, &requestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashPass, err := utils.HashPassword(requestDto.Password)
	if err != nil {
		fmt.Errorf("Error hash password", err)
	}

	requestDto.Password = hashPass

	h.service.Register(r.Context(), &requestDto)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID parameters", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserById(r.Context(), &id)

	if err != nil {
		fmt.Println("ошибка получения пользователя")
		return
	}

	fmt.Println(user)
}

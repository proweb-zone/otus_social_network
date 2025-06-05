package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/app/service"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
	"otus_social_network/app/internal/utils"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

type Handler struct {
	service *service.UserService
}

func Init(config *config.Config) *Handler {

	masterURL := []string{config.UrlsDb.DbMaster}
	slaveURLs := []string{
		config.UrlsDb.DbSlave1,
		config.UrlsDb.DbSlave2,
		config.UrlsDb.DbSlave3,
	}

	dataSource, err := postgres.NewReplicationRoutingDataSource(masterURL, slaveURLs, true)
	if err != nil {
		log.Fatal(err)
	}
	//defer dataSource.Close()

	userRepository := repository.InitPostgresRepository(dataSource)
	service := service.InitUserService(userRepository)

	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestDto dto.AuthRequestDto
	if err := utils.DecodeJson(body, &requestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestDto.Password = strings.TrimSpace(requestDto.Password)
	if len(requestDto.Password) == 0 {
		http.Error(w, "Error: field Password not found", http.StatusBadRequest)
		return
	}

	requestDto.Email = strings.TrimSpace(requestDto.Email)
	if len(requestDto.Email) == 0 {
		http.Error(w, "Error: field Email not found", http.StatusBadRequest)
		return
	}

	isValidEmail := utils.IsValidEmail(requestDto.Email)

	if !isValidEmail {
		http.Error(w, "Error: field Email invalid", http.StatusBadRequest)
		return
	}

	authResponse, err := h.service.Login(r.Context(), &requestDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(authResponse, w)
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

	requestDto.Password = strings.TrimSpace(requestDto.Password)
	if len(requestDto.Password) == 0 {
		http.Error(w, "Error: field password not found", http.StatusBadRequest)
		return
	}

	requestDto.Email = strings.TrimSpace(requestDto.Email)
	if len(requestDto.Email) == 0 {
		http.Error(w, "Error: field Email not found", http.StatusBadRequest)
		return
	}

	requestDto.First_name = strings.TrimSpace(requestDto.First_name)
	if len(requestDto.First_name) == 0 {
		http.Error(w, "Error: field First_name not found", http.StatusBadRequest)
		return
	}

	hashPass, err := utils.HashPassword(requestDto.Password)
	if err != nil {
		fmt.Errorf("Error: hash password", err)
		return
	}

	requestDto.Password = hashPass

	userResponse, err := h.service.Register(r.Context(), &requestDto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(userResponse, w)
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Error: invalid ID parameter", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserById(r.Context(), &id)

	if err != nil {
		http.Error(w, "Error: user not found", http.StatusBadRequest)
		return
	}

	utils.ResponseJson(user, w)
}

func (h *Handler) SearchUser(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	query := chi.URLParam(r, "query")
	prepairQuery := strings.Split(query, " ")

	if len(prepairQuery) < 2 {
		http.Error(w, "Error: first_name or last_name not found", http.StatusBadRequest)
		return
	}

	firstName := prepairQuery[0]
	lastName := prepairQuery[1]

	users, err := h.service.SearchUser(firstName, lastName)

	if err != nil {
		http.Error(w, "Error: users not found", http.StatusBadRequest)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf(" резултат выполнения за %s\n", elapsed)

	utils.ResponseJson(users, w)
}

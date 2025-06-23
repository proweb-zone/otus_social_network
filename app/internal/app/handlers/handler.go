package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
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
	userService    *service.UserService
	friendsService *service.FriendsService
	postService    *service.PostService
}

func Init(config *config.Config) *Handler {

	masterURL := []string{config.UrlsDb.DbMaster}
	slaveURLs := []string{
		config.UrlsDb.DbMaster,
		config.UrlsDb.DbMaster,
	}

	dataSource, err := postgres.NewReplicationRoutingDataSource(masterURL, slaveURLs, true)
	if err != nil {
		log.Fatal(err)
	}
	//defer dataSource.Close()

	userRepository := repository.InitUserRepository(dataSource)
	userService := service.InitUserService(userRepository)

	friendsRepository := repository.InitFriendsRepository(dataSource)
	friendsService := service.InitFriendsService(friendsRepository)

	postsRepository := repository.InitPostRepository(dataSource)
	postsService := service.InitPostService(postsRepository)

	return &Handler{
		userService:    userService,
		friendsService: friendsService,
		postService:    postsService,
	}
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

	authResponse, err := h.userService.Login(r.Context(), &requestDto)
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

	userResponse, err := h.userService.Register(r.Context(), &requestDto)
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

	user, err := h.userService.GetUserById(r.Context(), id)

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

	users, err := h.userService.SearchUser(firstName, lastName)

	if err != nil {
		http.Error(w, "Error: users not found", http.StatusBadRequest)
		return
	}

	elapsed := time.Since(start)
	fmt.Printf(" резултат выполнения за %s\n", elapsed)

	utils.ResponseJson(users, w)
}

func (h *Handler) SetFriend(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	friendIdStr := chi.URLParam(r, "user_id")

	friendId, err := strconv.Atoi(friendIdStr)
	if err != nil {
		http.Error(w, "Error: invalid USER ID parameter", http.StatusBadRequest)
		return
	}

	if userId == friendId {
		http.Error(w, "Error нельзя добавить себя в друзья", http.StatusBadRequest)
		return
	}

	// проверяем есть ли вообще такой пользователь в БД
	_, errService := h.userService.GetUserById(r.Context(), friendId)
	if errService != nil {
		http.Error(w, "Пользователь с таким id "+friendIdStr+" не существует", http.StatusBadRequest)
		return
	}

	// проверяем есть ли пользователь в друзья
	friend, _ := h.friendsService.GetFriendById(userId, friendId)
	if friend != nil {
		http.Error(w, "Такой пользователь уже в друзьях", http.StatusBadRequest)
		return
	}

	statusSetFriend, errSetFriend := h.friendsService.SetFriend(userId, friendId)
	if errSetFriend != nil {
		http.Error(w, "Ошибка добавления в друзья", http.StatusBadRequest)
		return
	}

	if statusSetFriend == "success" {
		w.Write([]byte("Успешное добавление пользователя в друзья!"))
		return
	}

}

func (h *Handler) DeleteFriend(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	friendIdStr := chi.URLParam(r, "user_id")
	friendId, err := strconv.Atoi(friendIdStr)
	if err != nil {
		http.Error(w, "Error: invalid USER ID parameter", http.StatusBadRequest)
		return
	}

	// проверяем есть ли пользователь в друзья
	_, errCheckFriend := h.friendsService.GetFriendById(userId, friendId)
	if errCheckFriend != nil {
		http.Error(w, "Пользователя нет в друзьях", http.StatusBadRequest)
		return
	}

	_, errDeleteFriend := h.friendsService.DeleteFriend(userId, friendId)
	if errDeleteFriend != nil {
		http.Error(w, "Ошибка удаления друга", http.StatusBadRequest)
		return
	}

	w.Write([]byte("Пользователь успешно удален"))
	return
}

func (h *Handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestDto dto.PostRequestDto
	if err := utils.DecodeJson(body, &requestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestDto.User_id = userId

	errCreatePost := h.postService.CreatePost(r.Context(), &requestDto)
	if errCreatePost != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// var postResponse dto.PostResponseDto
	// postResponse.Post_id = newPost.ID

	// utils.ResponseJson(postResponse, w)
}

func (h *Handler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestDto dto.PostRequestDto
	if err := utils.DecodeJson(body, &requestDto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestDto.User_id = userId

	errorPost := h.postService.UpdatePost(r.Context(), &requestDto)
	if errorPost != nil {
		http.Error(w, errorPost.Error(), http.StatusBadRequest)
		return
	}

	post, errPost := h.postService.GetPost(requestDto.User_id, requestDto.ID)
	if errPost != nil {
		http.Error(w, errPost.Error(), http.StatusBadRequest)
		return
	}

	utils.ResponseJson(post, w)

	w.Write([]byte("success"))
}

func (h *Handler) DeletePost(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	postIdStr := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Error: invalid USER ID parameter", http.StatusBadRequest)
		return
	}

	errorPost := h.postService.DeletePost(userId, postId)
	if errorPost != nil {
		http.Error(w, "Error: delete post", http.StatusBadRequest)
		return
	}

	w.Write([]byte("success"))
}

func (h *Handler) GetPost(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	postIdStr := chi.URLParam(r, "id")
	postId, err := strconv.Atoi(postIdStr)
	if err != nil {
		http.Error(w, "Error: invalid USER ID parameter", http.StatusBadRequest)
		return
	}

	post, errPost := h.postService.GetPost(userId, postId)
	if errPost != nil {
		http.Error(w, "Error: Пост с id "+postIdStr+"  не найден", http.StatusBadRequest)
		return
	}

	utils.ResponseJson(post, w)
}

func (h *Handler) FeedPost(w http.ResponseWriter, r *http.Request) {
	auth, errAccessToken := h.checkTokenAccess(r)

	if errAccessToken != nil {
		http.Error(w, "Error check Bearer Token", http.StatusBadRequest)
		return
	}

	userId := auth.User_id

	ids, errIds := h.friendsService.GetFriendIds(userId)
	if errIds != nil {
		http.Error(w, errIds.Error(), http.StatusBadRequest)
		return
	}

	posts, errPosts := h.postService.FeedPost(ids)
	if errPosts != nil {
		http.Error(w, errIds.Error(), http.StatusBadRequest)
		return
	}

	if posts == nil {
		http.Error(w, "news list friends not found", http.StatusBadRequest)
		return
	}

	utils.ResponseJson(posts, w)
}

func (h *Handler) checkTokenAccess(r *http.Request) (*entity.Auth, error) {
	// Извлечение токена из заголовка Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("Authorization header missing")
	}

	bearerToken := strings.Split(authHeader, " ")
	if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
		return nil, fmt.Errorf("Invalid authorization header format")
	}

	token := bearerToken[1]

	auth, err := h.userService.CheckAccessToken(token)
	if err != nil {
		return nil, err
	}

	return auth, nil
}

package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"otus_social_network/app/internal/app/handlers"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/app/service"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	eventclient "github.com/proweb-zone/event-client"
)

type HTTPServer struct {
	handlers *handlers.Handler
}

func StartServer(config *config.Config) {

	// init service dialog
	conn := postgres.Connect(config)
	defer postgres.Close(conn)

	masterURL := []string{config.UrlsDb.DbMaster}
	slaveURLs := []string{
		config.UrlsDb.DbMaster,
		config.UrlsDb.DbMaster,
	}

	dataSource, err := postgres.NewReplicationRoutingDataSource(masterURL, slaveURLs, true)
	if err != nil {
		log.Fatal(err)
	}

	// connect event client
	client, err := eventclient.New(eventclient.Config{
		GatewayAddress: config.GrpcServer.Addr,
		ServiceName:    "otus-service",
		MaxRetries:     5,
		RetryDelay:     1 * time.Second,
	})

	if err != nil {
		log.Fatalf("Failed to create event client: %v", err)
	}

	defer client.Close()

	log.Println("MS otus service started")

	userRepository := repository.InitUserRepository(dataSource)
	userService := service.InitUserService(userRepository, client)

	friendsRepository := repository.InitFriendsRepository(dataSource)
	friendsService := service.InitFriendsService(friendsRepository)

	postsRepository := repository.InitPostRepository(dataSource)
	postsService := service.InitPostService(postsRepository)

	dialogService := service.NewDialogService(friendsRepository, client)

	handlers := handlers.Init(config, userService, friendsService, postsService, dialogService)

	// subscribe on event
	err = client.Subscribe(context.Background(), []string{
		"dialog.send",
	}, handlers.Events)

	if err != nil {
		log.Fatalf("Failed to subscribe to events: %v", err)
	}

	r := chi.NewRouter()
	r.Post("/login", handlers.Login)
	r.Post("/user/register", handlers.Register)
	r.Get("/user/search/{query}", handlers.SearchUser)
	r.Get("/user/get/{id}", handlers.GetUser)

	r.Put("/friend/set/{user_id}", handlers.SetFriend)
	r.Put("/friend/delete/{user_id}", handlers.DeleteFriend)

	r.Post("/post/create", handlers.CreatePost)
	r.Put("/post/update", handlers.UpdatePost)
	r.Put("/post/delete/{id}", handlers.DeletePost)
	r.Get("/post/get/{id}", handlers.GetPost)
	r.Get("/post/feed", handlers.FeedPost)

	r.Get("/post/feed/posted", handlers.FeedPostWebSocket)

	go func() {
		http.ListenAndServe(":"+config.HTTPServer.ServerPort, r)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Shutting down...")
}

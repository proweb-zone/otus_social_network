package app

import (
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/server"
)

func InitApp(config *config.Config) {
	server.StartServer(config)
}

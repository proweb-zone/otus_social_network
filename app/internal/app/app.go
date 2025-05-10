package app

import (
	"otus_social_network/internal/config"
	"otus_social_network/internal/server"
)

func InitApp(config *config.Config) {
	server.StartServer(config)
}

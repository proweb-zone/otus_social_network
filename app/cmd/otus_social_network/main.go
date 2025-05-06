package main

import (
	"otus_social_network/app/internal/app"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/utils"
	"otus_social_network/app/pkg/logger"
)

func main() {
	currentDir := utils.GetProjectPath()

	configPath := config.ParseConfigPathFromCl(currentDir)
	сfgEnv := config.MustInit(configPath)

	zerologLogger := logger.ConfigureLogger(сfgEnv.Env)
	logger := logger.NewLogger(zerologLogger)

	app.InitApp(сfgEnv, currentDir, logger)
}

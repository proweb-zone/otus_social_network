package main

import (
	"otus_social_network/internal/app"
	"otus_social_network/internal/config"
	"otus_social_network/internal/utils"
)

func main() {
	currentDir := utils.GetProjectPath()
	configPath := config.ParseConfigPathFromCl(currentDir)
	config := config.MustInit(configPath)
	app.InitApp(config)
}

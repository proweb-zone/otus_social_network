package main

import (
	"otus_social_network/app/internal/app"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/utils"
)

func main() {
	currentDir := utils.GetProjectPath()
	configPath := config.ParseConfigPathFromCl(currentDir)
	config := config.MustInit(configPath)
	app.InitApp(config)
}

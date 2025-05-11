package utils

import (
	"os"
	"path/filepath"
)

func GetProjectPath() string {
	projectPath := os.Getenv("PROJECT_PATH")

	if projectPath != "" {
		return projectPath
	}

	currentDir, _ := os.Getwd()

	// Создаем путь к родительской директории
	parentDir := filepath.Join(currentDir, "../../../")

	return parentDir
}

package utils

import (
	"os"
	"path/filepath"
)

// func DecodeJson(body []byte, result any) error {
// 	if err := json.Unmarshal(body, result); err != nil {
// 		return fmt.Errorf("ошибка при декодировании JSON: %w", err)
// 	}

// 	return nil
// }

// func ConvertConfigToMap(config config.Config) map[string]any {
// 	result := make(map[string]any)
// 	value := reflect.ValueOf(config)

// 	for i := range value.NumField() {
// 		field := value.Type().Field(i)
// 		fieldValue := value.Field(i)

// 		yamlTag := field.Tag.Get("yaml")
// 		if yamlTag != "" {
// 			result[yamlTag] = fieldValue.Interface()
// 		}
// 	}

// 	return result
// }

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

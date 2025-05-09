package utils

import (
	"os"
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
	return currentDir
}

// func GetProjectRoot() (string, error) {
// 	currentDir, err := os.Getwd()
// 	if err != nil {
// 		return "", err
// 	}

// 	for {
// 		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
// 			return currentDir, nil
// 		}
// 		parentDir := filepath.Dir(currentDir)
// 		if parentDir == currentDir {
// 			break
// 		}
// 		currentDir = parentDir
// 	}

// 	return "", fmt.Errorf("project root not found")
// }

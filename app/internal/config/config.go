package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	Db         `yaml:"db"`
}

type HTTPServer struct {
	Address    string `yaml:"address" env-default:"localhost:8080"`
	Timeout    int    `yaml:"timeout" env-default:"2"`
	ServerPort string `yaml:"server_port"`
}

type Db struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Name     string `yaml:"db_Name"`
	User     string `yaml:"db_User"`
	Password string `yaml:"db_Password"`
	Option   string `yaml:"db_option"`
	StrConn  string `yaml:"db_url"`
}

func MustInit(configPath string) *Config {
	godotenv.Load(configPath)

	return &Config{
		Env: MustGetEnv("ENV"),
		HTTPServer: HTTPServer{
			Address:    MustGetEnv("HTTP_SERVER_ADDRESS"),
			Timeout:    MustGetEnvAsInt("HTTP_TIMEOUT"),
			ServerPort: MustGetEnv("SERVER_PORT"),
		},
		Db: Db{
			Driver:   MustGetEnv("DB_DRIVER"),
			Host:     MustGetEnv("DB_HOST"),
			Port:     MustGetEnv("DB_PORT"),
			Name:     MustGetEnv("DB_NAME"),
			User:     MustGetEnv("DB_USER"),
			Password: MustGetEnv("DB_PASSWORD"),
			Option:   MustGetEnv("DB_OPTION"),
			StrConn:  MustGetEnv("DB_URL"),
		},
	}
}

func PathDefault(workDir string, filename *string) string {
	if filename == nil {
		return filepath.Join(workDir, ".env")
	}

	return filepath.Join(workDir, *filename)
}

func MustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("no variable in env: %s", key)
	}
	return value
}

func MustGetEnvAsInt(name string) int {
	valueStr := MustGetEnv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return -1
}

func ParseConfigPathFromCl(currentDir string) string {
	return PathDefault(currentDir, nil)
}

package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local"`
	HTTPServer `yaml:"http_server"`
	UrlsDb
}

type HTTPServer struct {
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
}

type UrlsDb struct {
	DbMaster string
	DbSlave1 string
	DbSlave2 string
}

func MustInit(configPath string) *Config {
	godotenv.Load(configPath)

	dbMaster := &Db{
		Driver:   MustGetEnv("DB_DRIVER_MASTER"),
		Host:     MustGetEnv("DB_HOST_MASTER"),
		Port:     MustGetEnv("DB_PORT_MASTER"),
		Name:     MustGetEnv("DB_NAME_MASTER"),
		User:     MustGetEnv("DB_USER_MASTER"),
		Password: MustGetEnv("DB_PASSWORD_MASTER"),
		Option:   MustGetEnv("DB_OPTION_MASTER"),
	}

	var urlDbMaster = buildDbConnectUrl(dbMaster)

	dbSlave1 := &Db{
		Driver:   MustGetEnv("DB_DRIVER_SLAVE_1"),
		Host:     MustGetEnv("DB_HOST_SLAVE_1"),
		Port:     MustGetEnv("DB_PORT_SLAVE_1"),
		Name:     MustGetEnv("DB_NAME_SLAVE_1"),
		User:     MustGetEnv("DB_USER_SLAVE_1"),
		Password: MustGetEnv("DB_PASSWORD_SLAVE_1"),
		Option:   MustGetEnv("DB_OPTION_SLAVE_1"),
	}

	var urlDbSlave1 = buildDbConnectUrl(dbSlave1)

	dbSlave2 := &Db{
		Driver:   MustGetEnv("DB_DRIVER_SLAVE_2"),
		Host:     MustGetEnv("DB_HOST_SLAVE_2"),
		Port:     MustGetEnv("DB_PORT_SLAVE_2"),
		Name:     MustGetEnv("DB_NAME_SLAVE_2"),
		User:     MustGetEnv("DB_USER_SLAVE_2"),
		Password: MustGetEnv("DB_PASSWORD_SLAVE_2"),
		Option:   MustGetEnv("DB_OPTION_SLAVE_2"),
	}

	var urlDbSlave2 = buildDbConnectUrl(dbSlave2)

	return &Config{
		Env: MustGetEnv("ENV"),
		HTTPServer: HTTPServer{
			ServerPort: MustGetEnv("SERVER_PORT"),
		},
		UrlsDb: UrlsDb{
			DbMaster: urlDbMaster,
			DbSlave1: urlDbSlave1,
			DbSlave2: urlDbSlave2,
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

func buildDbConnectUrl(db *Db) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		db.Driver,
		db.User,
		db.Password,
		db.Host,
		db.Port,
		db.Name,
		db.Option,
	)
}

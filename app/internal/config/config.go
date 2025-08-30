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
	GrpcServer `yaml:"grpc_server"`
	Db
	UrlsDb
}

type HTTPServer struct {
	ServerPort string `yaml:"server_port"`
}

type GrpcServer struct {
	Addr string `yaml:"grpc_server_address"`
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
}

func MustInit(configPath string) *Config {
	godotenv.Load(configPath)

	env := MustGetEnv("ENV")

	// определяем переменную db сервера в dev/prod режиме
	dbPort := "5432"
	dbHostMater := MustGetEnv("DB_HOST_MASTER")
	if env == "DEV" {
		dbHostMater = "localhost"
		dbPort = MustGetEnv("DB_PORT_MASTER")
	}

	dbMaster := &Db{
		Driver:   MustGetEnv("DB_DRIVER_MASTER"),
		Host:     dbHostMater,
		Port:     dbPort,
		Name:     MustGetEnv("DB_NAME_MASTER"),
		User:     MustGetEnv("DB_USER_MASTER"),
		Password: MustGetEnv("DB_PASSWORD_MASTER"),
		Option:   MustGetEnv("DB_OPTION_MASTER"),
	}

	var urlDbMaster = buildDbConnectUrl(dbMaster)

	// определяем переменную grpc сервера в dev/prod режиме
	grpcAddrServer := MustGetEnv("GRPC_SERVER_ADDRESS")
	if env == "DEV" {
		grpcAddrServer = MustGetEnv("GRPC_SERVER_ADDRESS_DEV")
	}

	return &Config{
		Env: MustGetEnv("ENV"),
		HTTPServer: HTTPServer{
			ServerPort: MustGetEnv("SERVER_PORT"),
		},
		Db: *dbMaster,
		GrpcServer: GrpcServer{
			Addr: grpcAddrServer,
		},
		UrlsDb: UrlsDb{
			DbMaster: urlDbMaster,
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

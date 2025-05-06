package app

import (
	"fmt"
	"net/http"

	"ms_baskets/internal/config"
	"ms_baskets/internal/db/postgres"
	"ms_baskets/internal/server"
	"ms_baskets/pkg/logger"

	"github.com/go-chi/chi"
)

// Инициализация
func InitApp(сfgEnv *config.Config, projectPath string, logger *logger.Logger) {
	const op = "app.InitApp"

	gormIns, err := postgres.Create(buildDbConnectUrl(сfgEnv))
	if err != nil {
		logger.Fatal(err, "Экземпляр ORM не был создан", nil)
	}

	err = postgres.EnableAutoMigrate(gormIns)
	if err != nil {
		logger.Fatal(err, "Ошибка создания миграции БД", nil)
	}

	err = postgres.Connect(gormIns)

	if err != nil {
		logger.Error(err, op, nil)
		return
	}

	defer func() {
		if err := postgres.Close(gormIns); err != nil {
			logger.Error(err, op, nil)
		}
	}()

	router := server.ConfigureRouting(chi.NewRouter(), gormIns, logger)
	Run(router, logger, сfgEnv)
}

// Запуск
func Run(router *chi.Mux, logger *logger.Logger, сfgEnv *config.Config) *http.Server {
	return server.StartServer(router, logger, сfgEnv)
}

func buildDbConnectUrl(сfgEnv *config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		сfgEnv.Db.Driver,
		сfgEnv.Db.User,
		сfgEnv.Db.Password,
		сfgEnv.Db.Host,
		сfgEnv.Db.Port,
		сfgEnv.Db.Name,
		сfgEnv.Db.Option,
	)
}

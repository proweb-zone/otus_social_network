package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"otus_social_network/internal/config"
	"otus_social_network/internal/db/postgres"
	"otus_social_network/internal/migrator"
	"otus_social_network/pkg/logger"
)

func main() {
	const op = "app.InitMigrate"

	currentDir, _ := os.Getwd()

	var configPath string
	flag.StringVar(&configPath, "config", config.PathDefault(currentDir, nil), "path to config file")

	var action string
	flag.StringVar(&action, "action", "up", "path to config file")
	flag.Parse()

	var сfgEnv *config.Config = config.MustInit(configPath)

	postgresIns, err := postgres.Create(buildDbConnectUrl(сfgEnv))
	err = postgres.Connect(postgresIns)

	zerologLogger := logger.ConfigureLogger(сfgEnv.Env)
	logger := logger.NewLogger(zerologLogger)

	if err != nil {
		logger.Error(err, op, nil)
		return
	}

	sqlDb, err := postgresIns.DB()
	migrator := migrator.MustGetNewMigrator(сfgEnv.Name)
	switchAndExecMigrateAction(action, migrator, sqlDb, logger)
}

func buildDbConnectUrl(сfgEnv *config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		сfgEnv.Driver,
		сfgEnv.User,
		сfgEnv.Password,
		сfgEnv.Host,
		сfgEnv.Port,
		сfgEnv.Name,
		сfgEnv.Option,
	)
}

func switchAndExecMigrateAction(action string, migrator *migrator.Migrator, conn *sql.DB, logger *logger.Logger) {
	switch action {
	case "up":
		if err := migrator.Up(conn); err != nil {
			logger.Fatal(err, "error migration:", nil)
		} else {
			logger.Info("migration success.", nil)
		}
	case "down":
		if err := migrator.Down(conn); err != nil {
			logger.Fatal(err, "error rollback migration:", nil)
		} else {
			logger.Info("migration success.", nil)
		}
	default:
		logger.Warn("indefinite action on migration:"+action+". Available 'up' or 'down'.", nil, nil)
	}
}

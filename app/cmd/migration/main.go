package main

import (
	"database/sql"
	"flag"
	"fmt"

	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
	"otus_social_network/app/internal/migrator"
	"otus_social_network/app/internal/utils"
)

func main() {

	currentDir := utils.GetProjectPath()
	configPath := config.PathDefault(currentDir, nil)
	config := config.MustInit(configPath)

	var action string
	flag.StringVar(&action, "action", "up", "path to config file")
	flag.Parse()

	sqlDb := postgres.Connect(config.Db.StrConn)
	defer postgres.Close(sqlDb)

	migrator := migrator.MustGetNewMigrator(config.Db.Name)
	switchAndExecMigrateAction(action, migrator, sqlDb)
}

func switchAndExecMigrateAction(action string, migrator *migrator.Migrator, conn *sql.DB) {
	switch action {
	case "up":
		if err := migrator.Up(conn); err != nil {
			fmt.Println("Up migration failed")
		} else {
			fmt.Println("Up migration success.")
		}
	case "down":
		if err := migrator.Down(conn); err != nil {
			fmt.Println("Down migration failed")
		} else {
			fmt.Println("Down migration success.")
		}
	default:
		//logger.Warn("indefinite action on migration:"+action+". Available 'up' or 'down'.", nil, nil)
	}
}

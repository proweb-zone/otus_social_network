package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"

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

	// sqlDb := postgres.Connect(config)
	// defer postgres.Close(sqlDb)

	masterURL := []string{config.UrlsDb.DbMaster}
	slaveURLs := []string{
		config.UrlsDb.DbMaster,
	}

	sqlDb, err := postgres.NewReplicationRoutingDataSource(masterURL, slaveURLs, true)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	masterDb, err := sqlDb.GetDBMaster(ctx)
	if err != nil {
		fmt.Errorf("")
	}
	//defer sqlDb.Close()

	migrator := migrator.MustGetNewMigrator(config.Db.Name)
	switchAndExecMigrateAction(action, migrator, masterDb)
}

func switchAndExecMigrateAction(action string, migrator *migrator.Migrator, conn *sql.DB) {
	switch action {
	case "up":
		if err := migrator.Up(conn); err != nil {
			fmt.Println("Up migration failed", err)
		} else {
			fmt.Println("Up migration success.")
		}
	case "down":
		if err := migrator.Down(conn); err != nil {
			fmt.Println("Down migration failed", err)
		} else {
			fmt.Println("Down migration success.")
		}
	default:
		fmt.Println("parameter -action up or down not found")
	}
}

package postgres

import (
	"database/sql"
	"fmt"
	"otus_social_network/app/internal/config"
)

func Connect(config *config.Config) *sql.DB {
	connStr := buildDbConnectUrl(config)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func Close(db *sql.DB) error {
	return db.Close()
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

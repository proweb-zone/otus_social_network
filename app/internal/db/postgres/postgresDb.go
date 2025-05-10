package postgres

import "database/sql"

func Connect(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	return db
}

func Close(db *sql.DB) error {
	return db.Close()
}

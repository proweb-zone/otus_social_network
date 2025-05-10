package repository

import "database/sql"

type UserRepository struct {
	db *sql.DB
}

func InitPostgresRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func Login() {

}

func Register() {

}

func GetUserById() {

}

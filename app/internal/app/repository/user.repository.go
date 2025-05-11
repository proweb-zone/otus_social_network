package repository

import (
	"context"
	"database/sql"
	"otus_social_network/app/internal/app/entity"
)

type UserRepository struct {
	db *sql.DB
}

func InitPostgresRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Login() {

}

func (r *UserRepository) Create(ctx context.Context, user *entity.Users) (int, error) {
	stmt := `INSERT INTO users (first_name, last_name, email, password, birth_date, gender, hobby, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := r.db.Exec(stmt, user.First_name, user.Last_name, user.Email, user.Password, user.Birth_date, user.Gender, user.Hobby, user.City)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id *int) (*entity.Users, error) {
	row := r.db.QueryRow("SELECT id, first_name, last_name, email, birth_date, gender, hobby, city FROM users WHERE id = $1", id)

	var user entity.Users
	err := row.Scan(&user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Birth_date, &user.Gender, &user.Hobby, &user.City)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
)

type UserRepository struct {
	db *sql.DB
}

func InitPostgresRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.Users) (*dto.UsersResponseDto, error) {
	var userResponse dto.UsersResponseDto

	stmt := `INSERT INTO users (first_name, last_name, email, password, birth_date, gender, hobby, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.Exec(stmt, user.First_name, user.Last_name, user.Email, user.Password, user.Birth_date, user.Gender, user.Hobby, user.City)
	if err != nil {
		return &userResponse, err
	}

	return &userResponse, nil
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

func (r *UserRepository) GetUserByEmail(ctx context.Context, email *string) (*entity.Users, error) {
	row := r.db.QueryRow("SELECT id, first_name, last_name, email, password, birth_date, gender, hobby, city FROM users WHERE email = $1", email)

	var user entity.Users
	err := row.Scan(&user.ID, &user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Birth_date, &user.Gender, &user.Hobby, &user.City)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateToken(ctx context.Context, user *entity.Users, token *string) (*entity.Auth, error) {
	var auth entity.Auth

	stmt := `INSERT INTO auth (user_id, token) VALUES ($1, $2)`

	_, err := r.db.Exec(stmt, user.ID, token)
	if err != nil {
		return nil, fmt.Errorf("Error create token in DB")
	}

	row := r.db.QueryRow("SELECT token FROM auth WHERE token = $1", token)

	errToken := row.Scan(&auth.Token)
	if errToken != nil {
		return nil, fmt.Errorf("Error get token in DB")
	}

	return &auth, nil
}

func (r *UserRepository) GetTokenByUserId(ctx *context.Context, userId *uint) (*entity.Auth, error) {
	row := r.db.QueryRow("SELECT id, user_id, token, created_at FROM auth WHERE user_id = $1", userId)

	var auth entity.Auth
	err := row.Scan(&auth.ID, &auth.User_id, &auth.Token, &auth.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

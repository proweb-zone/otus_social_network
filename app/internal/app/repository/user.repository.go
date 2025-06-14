package repository

import (
	"context"
	"fmt"
	"log"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/db/postgres"
	"strings"
)

type UserRepository struct {
	dataSource *postgres.ReplicationRoutingDataSource
}

func InitPostgresRepository(dataSource *postgres.ReplicationRoutingDataSource) *UserRepository {
	return &UserRepository{dataSource}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.Users) (*entity.Users, error) {

	user.First_name = strings.ToLower(user.First_name)
	user.Last_name = strings.ToLower(user.Last_name)

	masterDb, err := r.dataSource.GetDBMaster(context.Background())
	if err != nil {
		return nil, err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer tx.Rollback()

	const insertQuery = `INSERT INTO users (first_name, last_name, email, password, birth_date, gender, hobby, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	stmt, err := masterDb.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	_, errExec := stmt.ExecContext(ctx, user.First_name, user.Last_name, user.Email, user.Password, user.Birth_date, user.Gender, user.Hobby, user.City)
	if errExec != nil {
		fmt.Errorf("Error ExecContext ", errExec)
		return nil, errExec
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	row := masterDb.QueryRow("SELECT id FROM users WHERE email = $1", user.Email)

	var newUser entity.Users
	err = row.Scan(&newUser.ID)

	if err != nil {
		return nil, err
	}

	fmt.Println(newUser.ID)
	return &newUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id *int) (*entity.Users, error) {
	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT id, first_name, last_name, email, birth_date, gender, hobby, city FROM users WHERE id = $1", id)

	var user entity.Users
	err := row.Scan(&user.ID, &user.First_name, &user.Last_name, &user.Email, &user.Birth_date, &user.Gender, &user.Hobby, &user.City)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email *string) (*entity.Users, error) {
	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT id, first_name, last_name, email, password, birth_date, gender, hobby, city FROM users WHERE email = $1", email)

	var user entity.Users
	err := row.Scan(&user.ID, &user.First_name, &user.Last_name, &user.Email, &user.Password, &user.Birth_date, &user.Gender, &user.Hobby, &user.City)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) CreateToken(ctx context.Context, user *entity.Users, token *string) (*entity.Auth, error) {
	var auth entity.Auth

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return &auth, err
	}

	stmt := `INSERT INTO auth (user_id, token) VALUES ($1, $2)`

	_, errQuery := masterDb.Exec(stmt, user.ID, token)
	if errQuery != nil {
		return nil, fmt.Errorf("Error create token in DB")
	}

	row := masterDb.QueryRow("SELECT token FROM auth WHERE token = $1", token)

	errToken := row.Scan(&auth.Token)
	if errToken != nil {
		return nil, fmt.Errorf("Error get token in DB")
	}

	return &auth, nil
}

func (r *UserRepository) GetTokenByUserId(ctx *context.Context, userId *uint) (*entity.Auth, error) {
	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT id, user_id, token, created_at FROM auth WHERE user_id = $1", userId)

	var auth entity.Auth
	err := row.Scan(&auth.ID, &auth.User_id, &auth.Token, &auth.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *UserRepository) CheckToken(token string) (*entity.Auth, error) {
	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT id, user_id, token, created_at FROM auth WHERE token = $1", token)

	var auth entity.Auth
	err := row.Scan(&auth.ID, &auth.User_id, &auth.Token, &auth.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &auth, nil
}

func (r *UserRepository) BatchInsertUsers(users []*entity.Users) error {

	ctx := context.Background()

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		fmt.Errorf("Error open transaction")
		log.Fatal(err)
	}
	defer tx.Rollback()

	const insertQuery = `INSERT INTO users (first_name, last_name, email, password, birth_date, gender, hobby, city) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	stmt, err := masterDb.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, user := range users {

		user.First_name = strings.ToLower(user.First_name)
		user.Last_name = strings.ToLower(user.Last_name)

		_, err := stmt.ExecContext(ctx, user.First_name, user.Last_name, user.Email, user.Password, user.Birth_date, user.Gender, user.Hobby, user.City)
		if err != nil {
			fmt.Errorf("Error ExecContext ", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Errorf("ошибка проведения транзацикции: %w", err)
		log.Fatal(err)
	}

	return nil
}

func (r *UserRepository) SearchUsers(firstName string, lastName string) ([]*entity.Users, error) {
	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	firstName = strings.ToLower(firstName)
	lastName = strings.ToLower(lastName)

	stmt, errPrepair := slaveDb.Prepare("SELECT id, first_name, last_name FROM users WHERE first_name LIKE $1 AND last_name LIKE $2 ORDER BY id")
	if errPrepair != nil {
		return nil, fmt.Errorf("prepare statement: %w", errPrepair)
	}
	defer stmt.Close()

	rows, err := stmt.Query("%"+firstName+"%", "%"+lastName+"%")
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}
	defer rows.Close()

	var users []*entity.Users
	for rows.Next() {
		user := new(entity.Users)
		err := rows.Scan(&user.ID, &user.First_name, &user.Last_name)
		if err != nil {
			return nil, fmt.Errorf("scan row: %w", err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return users, nil
}

package repository

import (
	"context"
	"fmt"
	"log"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/db/postgres"
)

type FriendsRepository struct {
	dataSource *postgres.ReplicationRoutingDataSource
}

func InitFriendsRepository(dataSource *postgres.ReplicationRoutingDataSource) *FriendsRepository {
	return &FriendsRepository{dataSource}
}

func (r *FriendsRepository) SetFriend(userId int, friendId int) (string, error) {

	ctx := context.Background()

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return "", err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer tx.Rollback()

	const insertQuery = `INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)`
	stmt, err := masterDb.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return "", err
	}
	defer stmt.Close()

	_, errExec := stmt.ExecContext(ctx, userId, friendId)
	if errExec != nil {
		fmt.Errorf("Error ExecContext ", errExec)
		return "", errExec
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return "success", nil
}

func (r *FriendsRepository) DeleteFriend(userId int, friendId int) (string, error) {

	ctx := context.Background()

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return "", err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer tx.Rollback()

	_, err = masterDb.Exec("DELETE FROM friends WHERE user_id = $1 AND friend_id = $2", userId, friendId)
	if err != nil {
		fmt.Println("Error deleting friend:", err)
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return "success", nil
}

func (r *FriendsRepository) GetFriendById(userId int, friendId int) (*entity.Friends, error) {

	slaveDb := r.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT user_id, friend_id FROM friends WHERE user_id = $1 AND friend_id = $2", userId, friendId)

	var friend entity.Friends
	err := row.Scan(&friend.User_id, &friend.Friend_id)

	if err != nil {
		return nil, err
	}

	return &friend, nil
}

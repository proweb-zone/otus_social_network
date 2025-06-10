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

func (r *FriendsRepository) SetFriend(userId int, friendId int) (*entity.Friends, error) {

	ctx := context.Background()

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer tx.Rollback()

	const insertQuery = `INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)`
	stmt, err := masterDb.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	_, errExec := stmt.ExecContext(ctx, userId, friendId)
	if errExec != nil {
		fmt.Errorf("Error ExecContext ", errExec)
		return nil, errExec
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return nil, nil
}

func (r *FriendsRepository) DeleteFriend(friendId int) (*entity.Friends, error) {

	ctx := context.Background()

	masterDb, err := r.dataSource.GetDBMaster(ctx)
	if err != nil {
		return nil, err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer tx.Rollback()

	const insertQuery = `"DELETE FROM friends WHERE friend_id = $1"`
	stmt, err := masterDb.PrepareContext(ctx, insertQuery)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer stmt.Close()

	_, errExec := stmt.ExecContext(ctx, friendId)
	if errExec != nil {
		return nil, fmt.Errorf("Error ExecContext ", errExec)
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return nil, nil
}

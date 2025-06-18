package repository

import (
	"context"
	"fmt"
	"log"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/db/postgres"
)

type PostRepository struct {
	dataSource *postgres.ReplicationRoutingDataSource
}

func InitPostRepository(dataSource *postgres.ReplicationRoutingDataSource) *PostRepository {
	return &PostRepository{dataSource}
}

func (p *PostRepository) GetPostById(userId int, postId int) (*entity.Posts, error) {
	slaveDb := p.dataSource.ChooseSlave()

	if slaveDb == nil {
		return nil, fmt.Errorf("no available slave databases")
	}

	row := slaveDb.QueryRow("SELECT id, user_id, text, created_at, updated_at FROM posts WHERE id = $1 AND user_id = $2", postId, userId)

	var post entity.Posts
	err := row.Scan(&post.ID, &post.User_id, &post.Text, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *PostRepository) CreatePost(userId int, postId int) (*entity.Posts, error) {
	return nil, nil
}

func (p *PostRepository) DeletePost(userId int, postId int) (string, error) {
	ctx := context.Background()

	masterDb, err := p.dataSource.GetDBMaster(ctx)
	if err != nil {
		return "", err
	}

	tx, err := masterDb.BeginTx(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer tx.Rollback()

	_, err = masterDb.Exec("DELETE FROM posts WHERE id = $1 AND user_id = $2", postId, userId)
	if err != nil {
		fmt.Println("Error deleting post:", err)
		return "", err
	}

	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return "success", nil
}

func (p *PostRepository) UpdatePost(userId int) (*entity.Posts, error) {
	return nil, nil
}

func (p *PostRepository) FeedPost(userId int) (*entity.Posts, error) {
	return nil, nil
}

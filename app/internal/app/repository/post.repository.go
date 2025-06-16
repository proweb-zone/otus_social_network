package repository

import (
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/db/postgres"
)

type PostRepository struct {
	dataSource *postgres.ReplicationRoutingDataSource
}

func InitPostRepository(dataSource *postgres.ReplicationRoutingDataSource) *PostRepository {
	return &PostRepository{dataSource}
}

func (p *PostRepository) GetPostByUserId(userId int) (*entity.Posts, error) {
	return nil, nil
}

func (p *PostRepository) CreatePost() (*entity.Posts, error) {
	return nil, nil
}

func (p *PostRepository) DeletePost(idPost int) (string, error) {
	return "", nil
}

func (p *PostRepository) UpdatePost() (*entity.Posts, error) {
	return nil, nil
}

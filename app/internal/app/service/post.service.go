package service

import (
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func InitPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(userId int) (*entity.Posts, error) {
	return nil, nil
}

func (p *PostService) UpdatePost(userId int, postId int) (*entity.Posts, error) {
	return nil, nil
}

func (p *PostService) DeletePost(userId int, postId int) (string, error) {
	return p.repo.DeletePost(userId, postId)
}

func (p *PostService) GetPost(userId int, postId int) (*entity.Posts, error) {
	return p.repo.GetPostById(userId, postId)
}

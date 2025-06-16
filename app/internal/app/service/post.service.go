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

func (p *PostService) CreatePost() (*entity.Posts, error) {
	return nil, nil
}

func (p *PostService) UpdatePost() (*entity.Posts, error) {
	return nil, nil
}

func (p *PostService) DeletePost() (string, error) {
	return "", nil
}

func (p *PostService) GetPost() (*entity.Posts, error) {
	return nil, nil
}

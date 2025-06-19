package service

import (
	"context"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func InitPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(ctx context.Context, request *dto.PostRequestDto) (*dto.PostResponseDto, error) {
	newPost, err := p.repo.CreatePost(
		ctx,
		&entity.Posts{
			User_id: request.User_id,
			Text:    request.Text,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error: Create post ", err)
	}

	var postResponse dto.PostResponseDto

	postResponse.Post_id = newPost.ID

	return &postResponse, nil
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

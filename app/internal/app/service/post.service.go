package service

import (
	"context"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/db/redis"
	"strings"
	"time"
)

type PostService struct {
	repo *repository.PostRepository
}

func InitPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(ctx context.Context, request *dto.PostRequestDto) error {
	_, errorCreatePost := p.repo.CreatePost(
		ctx,
		&entity.Posts{
			User_id: request.User_id,
			Text:    request.Text,
		},
	)

	if errorCreatePost != nil {
		return errorCreatePost
	}

	redis, errConnRedis := redis.InitRedisDb()
	if errConnRedis != nil {
		return errConnRedis
	}

	errAddMsg := redis.AddMsg(request.User_id, request.Text)

	if errAddMsg != nil {
		return errAddMsg
	}

	return nil
}

func (p *PostService) UpdatePost(ctx context.Context, request *dto.PostRequestDto) error {
	var updatedAt time.Time
	if len(request.Updated_at) > 0 {
		parsedTime, err := time.Parse("2006-01-02", strings.TrimSpace(request.Updated_at))
		if err != nil {
			return fmt.Errorf("Error: Incorect date in field updated_at")
		}
		updatedAt = parsedTime
	}

	return p.repo.UpdatePost(
		ctx,
		&entity.Posts{
			ID:        request.ID,
			User_id:   request.User_id,
			Text:      request.Text,
			UpdatedAt: updatedAt,
		},
	)
}

func (p *PostService) DeletePost(userId int, postId int) error {
	return p.repo.DeletePost(userId, postId)
}

func (p *PostService) GetPost(userId int, postId int) (*entity.Posts, error) {
	return p.repo.GetPostById(userId, postId)
}

func (p *PostService) FeedPost(ids []int) ([]*entity.Posts, error) {
	p.repo.FeedPost(ids)

	redis, errConnRedis := redis.InitRedisDb()
	if errConnRedis != nil {
		return nil, errConnRedis
	}

	posts, errGetRedisMsg := redis.GetMessages(ids)
	if errGetRedisMsg != nil {
		return nil, errGetRedisMsg
	}

	return posts, nil
}

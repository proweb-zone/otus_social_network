package service

import (
	"context"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
)

type FriendsService struct {
	repo *repository.FriendsRepository
}

func InitFriendsService(repo *repository.FriendsRepository) *FriendsService {
	return &FriendsService{repo: repo}
}

func (u *FriendsService) SetFriend(ctx *context.Context, userId int, friendId int) (*entity.Friends, error) {
	return u.repo.SetFriend(userId, friendId)
}

func (u *FriendsService) DeleteFriend(ctx *context.Context, friendId int) (*entity.Friends, error) {
	return u.repo.DeleteFriend()
}

package service

import (
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
)

type FriendsService struct {
	repo *repository.FriendsRepository
}

func InitFriendsService(repo *repository.FriendsRepository) *FriendsService {
	return &FriendsService{repo: repo}
}

func (u *FriendsService) GetFriendById(userId int, friendId int) (*entity.Friends, error) {
	return u.repo.GetFriendById(userId, friendId)
}

func (u *FriendsService) GetFriendIds(userId int) ([]int, error) {
	return u.repo.GetFriendIds(userId)
}

func (u *FriendsService) SetFriend(userId int, friendId int) (string, error) {
	return u.repo.SetFriend(userId, friendId)
}

func (u *FriendsService) DeleteFriend(userId int, friendId int) (string, error) {
	return u.repo.DeleteFriend(userId, friendId)
}

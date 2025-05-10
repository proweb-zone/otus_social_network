package service

import (
	"fmt"
	"otus_social_network/internal/app/repository"
)

type UserService struct {
	repo *repository.UserRepository
}

func InitUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) Login() {
	fmt.Println("call service Login")
}

func (u *UserService) Register() {
	fmt.Println("call service Register")
}

func (u *UserService) GetUserById() {
	fmt.Println("call service GetUserById")
}

package service

import (
	"context"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
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

func (u *UserService) Register(ctx context.Context, request *dto.UsersRequestDto) (int, error) {
	// err := request.Validate()
	// if err != nil {
	// 	return 0, err
	// }
	return u.repo.Create(
		ctx,
		&entity.Users{
			First_name: request.First_name,
			Last_name:  request.Last_name,
			Email:      request.Email,
			Password:   request.Password,
			Birth_date: request.Birth_date,
			Gender:     request.Gender,
			Hobby:      request.Hobby,
			City:       request.City,
		},
	)
}

func (u *UserService) GetUserById(ctx context.Context, id *int) (*entity.Users, error) {
	return u.repo.GetUserById(ctx, id)
}

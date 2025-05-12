package service

import (
	"context"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/utils"
)

type UserService struct {
	repo *repository.UserRepository
}

func InitUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (u *UserService) Login(ctx context.Context, requestDto *dto.AuthRequestDto) (*entity.Auth, error) {

	user, err := u.repo.GetUserByEmail(
		ctx,
		&requestDto.Email,
	)

	if err != nil {
		return nil, err
	}

	isValidPass, err := utils.CheckPassword(user.Password, requestDto.Password)

	if err != nil {
		return nil, err
	}

	if !isValidPass {
		return nil, fmt.Errorf("password invalid")
	}

	token := "dfdf"

	return u.repo.CreateToken(
		ctx,
		user,
		&token,
	)

}

func (u *UserService) Register(ctx context.Context, request *dto.UsersRequestDto) (int, error) {
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

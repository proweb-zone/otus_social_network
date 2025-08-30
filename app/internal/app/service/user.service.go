package service

import (
	"context"
	"fmt"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/utils"
	"strings"
	"time"

	eventclient "github.com/proweb-zone/event-client"
	pb "github.com/proweb-zone/event-client/gen/go"
)

type UserService struct {
	repo   *repository.UserRepository
	client *eventclient.EventClient
}

func InitUserService(repo *repository.UserRepository, newClient *eventclient.EventClient) *UserService {
	return &UserService{repo: repo, client: newClient}
}

func (u *UserService) Login(ctx context.Context, requestDto *dto.AuthRequestDto) (*dto.AuthResponseDto, error) {
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

	var authResponse dto.AuthResponseDto

	tokenByUserId, _ := u.repo.GetTokenByUserId(&ctx, &user.ID)

	if tokenByUserId != nil && len(tokenByUserId.Token) > 0 {
		authResponse.Bearer_token = tokenByUserId.Token
	} else {
		token := utils.GenerateToken(32)

		auth, err := u.repo.CreateToken(
			ctx,
			user,
			&token,
		)

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		authResponse.Bearer_token = auth.Token
	}

	u.client.Publish(context.Background(), &pb.Event{
		Type:    "user.auth",
		Source:  "otus-service",
		Payload: []byte(fmt.Sprintf(`{"User_id": %d, "Token": "%s"}`, user.ID, authResponse.Bearer_token)),
	})

	return &authResponse, nil
}

func (u *UserService) Register(ctx context.Context, request *dto.UsersRequestDto) (*dto.UsersResponseDto, error) {
	isvalidUser, _ := u.repo.GetUserByEmail(ctx, &request.Email)

	if isvalidUser != nil {
		return nil, fmt.Errorf("Error: User with this email has already been registered")
	}

	var birthTime time.Time
	if len(request.Birth_date) > 0 {
		parsedTime, err := time.Parse("2006-01-02", strings.TrimSpace(request.Birth_date))
		if err != nil {
			return nil, fmt.Errorf("Error: Incorect date in field birth_date")
		}
		birthTime = parsedTime
	}

	newUser, err := u.repo.Create(
		ctx,
		&entity.Users{
			First_name: request.First_name,
			Last_name:  request.Last_name,
			Email:      request.Email,
			Password:   request.Password,
			Birth_date: birthTime,
			Gender:     request.Gender,
			Hobby:      request.Hobby,
			City:       request.City,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Error: Create user ", err)
	}

	var userResponse dto.UsersResponseDto

	userResponse.User_id = newUser.ID

	return &userResponse, nil

}

func (u *UserService) GetUserById(ctx context.Context, id int) (*entity.Users, error) {
	return u.repo.GetUserById(ctx, id)
}

func (u *UserService) SearchUser(firstName string, lastName string) ([]*entity.Users, error) {
	return u.repo.SearchUsers(firstName, lastName)
}

func (u *UserService) CheckAccessToken(token string) (*entity.Auth, error) {
	return u.repo.CheckToken(token)
}

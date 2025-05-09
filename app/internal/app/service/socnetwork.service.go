package service

import (
	"otus_social_network/app/internal/app/repository"
)

type SocNetworkService struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *SocNetworkService {
	return &SocNetworkService{repo: repo}
}

// func (s *SocNetworkService) GetItem(ctx context.Context, id uint) (*entity.Baskets, error) {
// 	return s.repo.GetItem(ctx, id)
// }

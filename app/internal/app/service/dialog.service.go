package service

import (
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
)

type DialogService struct {
	repo *repository.DialogRepository
}

func InitDialogService(repo *repository.DialogRepository) *DialogService {
	return &DialogService{repo: repo}
}

func (d *DialogService) SendMsgUser(requestDialog *dto.DialogRequestDto) error {
	return d.repo.SendMsgUser(requestDialog)
}

func (d *DialogService) GetDialogList() ([]*entity.Dialog, error) {
	// TODO add getDialogList
	return nil, nil
}

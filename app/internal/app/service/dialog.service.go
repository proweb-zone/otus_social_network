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

func (d *DialogService) SendMsgUser(requestDialog *dto.DialogRequestDto) (*entity.Dialog, error) {
	return d.repo.SendMsgUser(&entity.Dialog{
		User_id_sender:    requestDialog.User_id_sender,
		User_id_recipient: requestDialog.User_id_recipient,
		Msg:               requestDialog.Msg,
	})
}

func (d *DialogService) GetDialogList(userIdSender int, userIdRecepient int) (*[]entity.Dialog, error) {
	return d.repo.GetDialogList(userIdSender, userIdRecepient)
}

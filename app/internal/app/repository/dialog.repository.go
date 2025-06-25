package repository

import (
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/db/postgres"
)

type DialogRepository struct {
	dataSource *postgres.ReplicationRoutingDataSource
}

func InitDialogRepository(dataSource *postgres.ReplicationRoutingDataSource) *DialogRepository {
	return &DialogRepository{dataSource}
}

func (d *DialogRepository) SendMsgUser(requestDialog *dto.DialogRequestDto) error {
	return nil
}

func (d *DialogRepository) GetDialogList() ([]*entity.Dialog, error) {
	return nil, nil
}

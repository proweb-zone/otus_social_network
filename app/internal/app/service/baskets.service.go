package service

import (
	"context"
	"ms_baskets/internal/app/dto"
	"ms_baskets/internal/app/entity"
	"ms_baskets/internal/app/repository"
	"ms_baskets/internal/utils"

	"golang.org/x/sync/errgroup"
)

type BasketsService struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *BasketsService {
	return &BasketsService{repo: repo}
}

func (s *BasketsService) GetItem(ctx context.Context, id uint) (*entity.Baskets, error) {
	return s.repo.GetItem(ctx, id)
}

func (s *BasketsService) GetList(ctx context.Context, filters []any, queryParamsRequest dto.QueryParamsDto) (*[]map[string]any, int64, error) {
	if err := queryParamsRequest.Validate(); err != nil {
		return nil, 0, err
	}

	fields := queryParamsRequest.ValidateAttrs(entity.GetBasketsAttrs())
	g, ctx := errgroup.WithContext(ctx)
	orderBy := utils.BuildOrderBy(queryParamsRequest.Sort, queryParamsRequest.Order, entity.GetBasketsAttrs())

	var count int64
	var items *[]map[string]any
	where := entity.Baskets{ProductId: queryParamsRequest.ProductId, UserId: queryParamsRequest.UserId}

	g.Go(func() error {
		var err error
		count, err = s.repo.CountListItems(ctx, &where)
		return err
	})

	g.Go(func() error {
		var err error
		items, err = s.repo.GetList(ctx, orderBy, filters, fields, &where)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, 0, err
	}

	return items, count, nil
}

func (s *BasketsService) Create(ctx context.Context, request *dto.BasketsRequestDto) (*uint, error) {
	err := request.Validate()
	if err != nil {
		return nil, err
	}
	return s.repo.Create(
		ctx,
		&entity.Baskets{
			UserId:    request.UserId,
			ProductId: request.ProductId,
			Quantity:  request.Quantity,
		},
	)
}

func (s *BasketsService) Update(id uint, request *dto.UpdateBasketsRequestDto, ctx context.Context) (int64, error) {
	err := request.Validate()
	if err != nil {
		return 0, err
	}
	baskets := entity.Baskets{
		ProductId: request.ProductId,
		Quantity:  request.Quantity,
	}

	return s.repo.Update(ctx, id, &baskets)
}

func (s *BasketsService) Delete(id uint, ctx context.Context) (int64, error) {
	return s.repo.Delete(ctx, id)
}

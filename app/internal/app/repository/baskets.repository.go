package repository

import (
	"context"

	"ms_baskets/internal/app/entity"
	"ms_baskets/internal/app/enum"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Repository struct {
	Collection *mongo.Collection
	gormIns    *gorm.DB
}

func InitPostgresRepository(gormIns *gorm.DB) *Repository {
	return &Repository{gormIns: gormIns}
}

func (r *Repository) GetItem(ctx context.Context, id uint) (*entity.Baskets, error) {
	var item entity.Baskets

	result := r.gormIns.WithContext(ctx).Model(&entity.Baskets{}).First(&item, &entity.Baskets{Id: id})
	err := result.Error

	return &item, err
}

func (r *Repository) GetList(ctx context.Context, orderBy string, filters []any, fields []string, where *entity.Baskets) (*[]map[string]any, error) {
	query := r.gormIns.WithContext(ctx).Model(&entity.Baskets{})
	if len(fields) == 0 {
		for key, _ := range entity.GetBasketsAttrs() {
			fields = append(fields, key)
		}
	}

	query.Select(fields)
	query, initialCapacity := r.applyPagination(query, filters)
	query = query.Order(orderBy)

	list := make([]map[string]any, 0, initialCapacity)
	err := query.Find(&list, *where).Error

	return &list, err
}

func (r *Repository) CountListItems(ctx context.Context, where *entity.Baskets) (int64, error) {
	var count int64

	err := r.gormIns.WithContext(ctx).Model(&entity.Baskets{}).Select("id").Where(*where).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, err
}

func (r *Repository) Create(ctx context.Context, baskets *entity.Baskets) (*uint, error) {
	result := r.gormIns.WithContext(ctx).Model(&entity.Baskets{}).Create(&baskets)

	return &baskets.Id, result.Error
}

func (r *Repository) Update(ctx context.Context, id uint, item *entity.Baskets) (int64, error) {
	result := r.gormIns.WithContext(ctx).Model(item).Where("id = ?", id).Updates(item)

	return result.RowsAffected, result.Error
}

func (r *Repository) Delete(ctx context.Context, id uint) (int64, error) {
	result := r.gormIns.WithContext(ctx).Delete(&entity.Baskets{Id: id})

	return result.RowsAffected, result.Error
}

func (r *Repository) applyPagination(query *gorm.DB, filters []any) (*gorm.DB, int) {
	initialCapacity := 100

	if offset, ok := filters[enum.Offset].(int); ok {
		query = query.Offset(offset)
	}

	if limit, ok := filters[enum.Limit].(int); ok && limit > 0 {
		initialCapacity = limit
		query = query.Limit(limit)
	} else {
		query = query.Limit(100)
	}

	return query, initialCapacity
}

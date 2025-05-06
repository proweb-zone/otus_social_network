package postgres

import (
	"context"
	"errors"
	"ms_baskets/internal/app/entity"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

func Connect(gormIns *gorm.DB) error {
	sqlDB, err := gormIns.DB()

	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	sqlDB.SetConnMaxIdleTime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return sqlDB.PingContext(ctx)
}

func Close(gormIns *gorm.DB) error {
	if gormIns == nil {
		return errors.New("экземпляр ORM не получен")
	}

	sqlDB, _ := gormIns.DB()

	return sqlDB.Close()
}

func Create(dsn string) (*gorm.DB, error) {
	return gorm.Open(
		postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}),
		&gorm.Config{
			Logger:                 gormLogger.Default.LogMode(gormLogger.Warn),
			SkipDefaultTransaction: true,
			PrepareStmt:            false,
		},
	)
}

func EnableAutoMigrate(gormIns *gorm.DB) error {
	return gormIns.AutoMigrate(&entity.Baskets{})
}

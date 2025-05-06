package integration

import (
	"context"
	"fmt"
	"ms_baskets/internal/config"
	"ms_baskets/internal/db/postgres"
	"ms_baskets/internal/utils"
	"os/exec"
	"path/filepath"
	"testing"

	"gorm.io/gorm"
)

type ErrorfFunc func(format string, args ...any)

func SetupTest(t *testing.T) (*gorm.DB, *config.Config) {
	configs, err := getConfig(".env.test", t.Errorf)
	if err != nil {
		t.Fatalf("Файл конфигурации приложения не найден: %v", err)
	}

	postgresIns := settingUpTestEnv(configs)
	err = postgres.EnableAutoMigrate(postgresIns)
	if err != nil {
		t.Fatalf("Ошибка создания миграции БД: %v", err)
	}
	if err := postgres.Connect(postgresIns); err != nil {
		t.Fatalf("Не удалось подключиться к postgres: %v", err)
	}

	return postgresIns, configs
}

func getConfig(configFileName string, errorf ErrorfFunc) (*config.Config, error) {
	root, err := utils.GetProjectRoot()
	if err != nil {
		errorf("Корневой каталог не найден %v", err)
		return nil, err
	}

	configPath := config.PathDefault(root, &configFileName)
	return config.MustInit(configPath), nil
}

func buildDbConnectUrl(сfgEnv *config.Config) string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s?%s",
		сfgEnv.Db.Driver,
		сfgEnv.Db.User,
		сfgEnv.Db.Password,
		сfgEnv.Db.Host,
		сfgEnv.Db.Port,
		сfgEnv.Db.Name,
		сfgEnv.Db.Option,
	)
}

func settingUpTestEnv(configs *config.Config) *gorm.DB {
	postgresIns, err := postgres.Create(buildDbConnectUrl(configs))
	if err != nil {
		fmt.Printf("Ошибка при подключении к postgres: %v", err)

		return nil
	}

	return postgresIns
}

func PgContainerUp() error {
	root, err := utils.GetProjectRoot()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker-compose", "-f", filepath.Join(root, "compose.test.yml"), "up", "-d")
	return cmd.Run()
}

func PgContainerDown() error {
	root, err := utils.GetProjectRoot()
	if err != nil {
		return err
	}

	cmd := exec.Command("docker-compose", "-f", filepath.Join(root, "compose.test.yml"), "down", "-v")
	return cmd.Run()
}

func TruncateTable(gormIns *gorm.DB, ctx context.Context) error {
	sqlDb, _ := gormIns.DB()
	_, err := sqlDb.ExecContext(ctx, `TRUNCATE TABLE baskets`)

	return err
}

func DropTable(gormIns *gorm.DB, ctx context.Context) error {
	sqlDb, _ := gormIns.DB()
	_, err := sqlDb.ExecContext(ctx, `DROP TABLE baskets`)

	return err
}

package baskets_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"ms_baskets/internal/app/controller"
	"ms_baskets/internal/app/dto"
	"ms_baskets/internal/app/entity"
	"ms_baskets/internal/config"
	"ms_baskets/internal/db/postgres"
	"ms_baskets/pkg/logger"
	"ms_baskets/test/integration"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BasketsTestSuite struct {
	suite.Suite
	gormIns    *gorm.DB
	configs    *config.Config
	controller *controller.Controller
	ctx        context.Context
	testBasket dto.BasketsRequestDto
}

type Result[T any] struct {
	Result T `json:"result"`
}

func (s *BasketsTestSuite) createBasket() (*uint, error) {
	newBasket := &entity.Baskets{
		UserId:    s.testBasket.UserId,
		ProductId: s.testBasket.ProductId,
		Quantity:  s.testBasket.Quantity,
	}

	result := s.gormIns.Model(&entity.Baskets{}).Create(&newBasket)
	err := result.Error

	if err != nil {
		return nil, fmt.Errorf("ошибка при добавлени записи: %w", err)
	}

	return &newBasket.Id, err
}

func (s *BasketsTestSuite) SetupSuite() {
	s.T().Log("Запуск тестового окружения")

	err := integration.PgContainerUp()
	assert.NoError(s.T(), err, "Ошибка при запуске контейнера postgres")
	time.Sleep(5 * time.Second) // время для инициализации БД

	s.ctx = context.Background()
	s.gormIns, s.configs = integration.SetupTest(s.T())

	zerologLogger := logger.ConfigureLogger(s.configs.Env)
	s.controller = controller.GetController(logger.NewLogger(zerologLogger), s.gormIns)

	s.testBasket = dto.BasketsRequestDto{
		UserId:    1,
		ProductId: 1,
		Quantity:  10,
	}
}

func (s *BasketsTestSuite) TearDownSuite() {
	s.T().Log("Удаление тестового окружения")

	err := integration.DropTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err, "Ошибка удаления таблицы")

	err = postgres.Close(s.gormIns)
	assert.NoError(s.T(), err, "Ошибка при закрытии соединения с postgres")

	err = integration.PgContainerDown()
	assert.NoError(s.T(), err, "Ошибка при удалении контейнера")
}

func (s *BasketsTestSuite) TestCreateBasketItem() {
	err := integration.TruncateTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err)

	body, err := json.Marshal(s.testBasket)
	assert.NoError(s.T(), err)

	req, err := http.NewRequest("PUT", "/api/baskets/create", bytes.NewBuffer(body))
	assert.NoError(s.T(), err)

	rr := httptest.NewRecorder()
	s.controller.Create(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code, "Неверный статус код")

	var response Result[dto.CreateResponse]

	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), *response.Result.Id, uint(0), "Неверный результат создания записи")
}

func (s *BasketsTestSuite) TestUpdateBasketItem() {
	err := integration.TruncateTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err)

	lastInsertId, err := s.createBasket()
	assert.NoError(s.T(), err)

	lastInsertIdStr := strconv.Itoa(int(*lastInsertId))

	body, err := json.Marshal(&dto.BasketsRequestDto{Quantity: 100, ProductId: 90})
	assert.NoError(s.T(), err)

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", lastInsertIdStr)

	req, err := http.NewRequest("POST", "/api/baskets/update/"+lastInsertIdStr, bytes.NewBuffer(body))
	assert.NoError(s.T(), err)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	rr := httptest.NewRecorder()
	s.controller.Update(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code, "Неверный статус код")

	var response Result[dto.UpdateResponse]

	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(s.T(), err)
	assert.Greater(s.T(), response.Result.Count, 0, "Неверный результат обновления записи")
}

func (s *BasketsTestSuite) TestDeleteBasketItem() {
	err := integration.TruncateTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err)

	lastInsertId, err := s.createBasket()
	assert.NoError(s.T(), err)

	lastInsertIdStr := strconv.Itoa(int(*lastInsertId))

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", lastInsertIdStr)

	req, err := http.NewRequest("DELETE", "/api/baskets/"+lastInsertIdStr, nil)
	assert.NoError(s.T(), err)

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	rr := httptest.NewRecorder()
	s.controller.Delete(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code, "Неверный статус код")

	var response Result[dto.DeleteResponse]

	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), int64(1), response.Result.Count, "Запись не была удалена")
}

func (s *BasketsTestSuite) TestGetBasketItem() {
	err := integration.TruncateTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err)

	lastInsertId, err := s.createBasket()
	assert.NoError(s.T(), err)

	lastInsertIdStr := strconv.Itoa(int(*lastInsertId))

	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("id", lastInsertIdStr)
	req, err := http.NewRequest("GET", "/api/baskets/"+lastInsertIdStr, nil)
	assert.NoError(s.T(), err)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	rr := httptest.NewRecorder()
	s.controller.GetItem(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code, "Неверный статус код")

	var response Result[dto.BasketsRequestDto]

	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(s.T(), err)

	expected := Result[dto.BasketsRequestDto]{Result: s.testBasket}

	assert.Equal(s.T(), expected, response, "Неверные данные в полученной записи")
}

func (s *BasketsTestSuite) TestGetBasketList() {
	err := integration.TruncateTable(s.gormIns, s.ctx)
	assert.NoError(s.T(), err)

	_, err = s.createBasket()
	assert.NoError(s.T(), err)

	routeCtx := chi.NewRouteContext()
	req, err := http.NewRequest("GET", "/api/baskets?limit=100&page=1&attrs[]=user_id&attrs[]=product_id", nil)
	assert.NoError(s.T(), err)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

	rr := httptest.NewRecorder()
	s.controller.GetList(rr, req)

	assert.Equal(s.T(), http.StatusOK, rr.Code, "Неверный статус код")

	var response Result[integration.TestItemsList[entity.Baskets]]

	err = json.NewDecoder(rr.Body).Decode(&response)
	assert.NoError(s.T(), err)

	expected := Result[integration.TestItemsList[entity.Baskets]]{Result: integration.TestItemsList[entity.Baskets]{
		Items: ToMap(&s.testBasket),
		Total: 1,
		Count: 1,
	}}

	assert.Equal(s.T(), expected, response, "Неверные данные в полученной записи")
}

func ToMap(b *dto.BasketsRequestDto) []entity.Baskets {
	return []entity.Baskets{
		{
			ProductId: b.ProductId,
			UserId:    b.UserId,
		},
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(BasketsTestSuite))
}

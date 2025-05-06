package controller

import (
	"errors"
	"io"
	"net/http"
	"strconv"

	"ms_baskets/internal/app/dto"
	"ms_baskets/internal/app/entity"
	"ms_baskets/internal/app/repository"
	"ms_baskets/internal/app/service"
	"ms_baskets/internal/utils"

	"ms_baskets/pkg/logger"
	"ms_baskets/pkg/respond"

	"github.com/go-chi/chi"
	"gorm.io/gorm"
)

type Result[T any] struct {
	Result *T `json:"result"`
}

type List[T any] struct {
	Count int   `json:"count"`
	Total int64 `json:"total"`
	Items *[]T  `json:"items"`
}

type Controller struct {
	logger  *logger.Logger
	service *service.BasketsService
}

func GetController(logger *logger.Logger, gormIns *gorm.DB) *Controller {
	basketsRepository := repository.InitPostgresRepository(gormIns)
	service := service.NewService(basketsRepository)

	return &Controller{logger: logger, service: service}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	body := c.parseBody(w, r)

	var request dto.BasketsRequestDto
	if err := utils.DecodeJson(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Получен некорректный формат JSON", err, request)
		return
	}

	resultId, err := c.service.Create(r.Context(), &request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "произошла ошибка при создании записи", errors.New(""), request)
		return
	}

	respond.Respond(w, r, http.StatusOK, Result[dto.CreateResponse]{Result: &dto.CreateResponse{Status: "success", Id: resultId}})
}

func (c *Controller) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Получен некорректный идентификатор", err, map[string]any{"id": id})
		return
	}
	body := c.parseBody(w, r)

	var request dto.UpdateBasketsRequestDto
	if err := utils.DecodeJson(body, &request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Получен некорректный формат JSON", err, request)
		return
	}

	count, err := c.service.Update(uint(idInt), &request, r.Context())
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Ошибка при обновлении записи", err, request)
		return
	}

	respond.Respond(w, r, http.StatusOK, Result[dto.UpdateResponse]{Result: &dto.UpdateResponse{Status: "success", Count: int(count)}})
}

func (c *Controller) GetItem(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Получен некорректный идентификатор", err, map[string]any{"id": id})
		return
	}
	result, err := c.service.GetItem(r.Context(), uint(idInt))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Ошибка получения записи", err, map[string]any{"id": id})
		return
	}

	respond.Respond(w, r, http.StatusOK, Result[entity.Baskets]{Result: result})
}

func (c *Controller) GetList(w http.ResponseWriter, r *http.Request) {
	paginationParams := dto.ParsePagination(r)
	queryParams := dto.ParseQueryParams(r)

	queryParamsRequest := dto.QueryParamsDto{
		UserId:    queryParams.UserId,
		ProductId: queryParams.ProductId,
		Sort:      queryParams.Sort,
		Order:     queryParams.Order,
		Attrs:     queryParams.Attrs,
	}

	filters := append(
		make([]any, 0, 2),
		int(paginationParams.Limit),
		paginationParams.Offset,
	)

	items, count, err := c.service.GetList(
		r.Context(),
		filters,
		queryParamsRequest,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Ошибка получения списка", err, nil)
		return
	}

	respond.Respond(w, r, http.StatusOK, Result[List[map[string]any]]{Result: &List[map[string]any]{Count: len(*items), Items: items, Total: count}})
}

func (c *Controller) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Получен некорректный идентификатор", err, map[string]any{"id": id})
		return
	}
	rowsAffected, err := c.service.Delete(uint(idInt), r.Context())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Ошибка удаления записи", err, map[string]any{"id": id})
		return
	}

	respond.Respond(w, r, http.StatusOK, Result[dto.DeleteResponse]{Result: &dto.DeleteResponse{Status: "success", Count: rowsAffected}})
}

func (c *Controller) handleError(w http.ResponseWriter, r *http.Request, msg string, err error, data any) {
	c.logger.Error(err, msg, map[string]any{
		"request": data,
		"error":   err.Error(),
	})

	respond.Respond(w, r, http.StatusBadRequest, dto.ErrorResponse{
		Error:  err.Error(),
		Result: msg,
	})
}

func (c *Controller) parseBody(w http.ResponseWriter, r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		c.handleError(w, r, "Ошибка чтения запроса", err, body)
		return nil
	}
	defer r.Body.Close()

	return body
}

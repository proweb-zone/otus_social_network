package dto

import (
	"net/http"
	"strconv"
)

type PaginationParams struct {
	Offset int
	Limit  int64
}

type QueryParams struct {
	ProductId uint
	UserId    uint
	Sort      string
	Order     string
	Attrs     []string
}

func ParsePagination(r *http.Request) PaginationParams {
	var offset int = 0
	var limit int64 = 100

	if r.URL.Query().Get("limit") != "" {
		if parsedLimit, err := strconv.ParseInt(r.URL.Query().Get("limit"), 10, 64); err == nil {
			limit = parsedLimit
		}
	}
	if r.URL.Query().Get("page") != "" {
		if num, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64); err == nil {
			offset = int(num-1) * int(limit)
		}
	}

	return PaginationParams{Offset: offset, Limit: limit}
}

func ParseQueryParams(r *http.Request) *QueryParams {
	q := r.URL.Query()
	var queryParams = QueryParams{Sort: "id", Order: "ASC"}

	if userId := q.Get("user_id"); userId != "" {
		if result, err := strconv.ParseUint(userId, 10, 64); err == nil {
			queryParams.UserId = uint(result)
		}
	}

	if productId := q.Get("product_id"); productId != "" {
		if result, err := strconv.ParseUint(productId, 10, 64); err == nil {
			queryParams.ProductId = uint(result)
		}
	}

	if sort := q.Get("sort"); sort != "" {
		queryParams.Sort = sort
	}

	if order := q.Get("order"); order != "" {
		queryParams.Order = order
	}

	if attrs := q["attrs[]"]; len(attrs) > 0 {
		queryParams.Attrs = attrs
	}

	return &queryParams
}

package integration

type TestItemsList[T any] struct {
	Count int `json:"count"`
	Total int `json:"total"`
	Items []T `json:"items"`
}

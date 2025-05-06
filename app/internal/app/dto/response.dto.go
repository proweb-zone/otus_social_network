package dto

type ErrorResponse struct {
	Result any    `json:"result"`
	Error  string `json:"error"`
}

type CreateResponse struct {
	Status string `json:"status"`
	Id     *uint  `json:"id"`
}

type UpdateResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type DeleteResponse struct {
	Status string `json:"status"`
	Count  int64  `json:"count"`
}

type ImageResponse struct {
	Image string `json:"image"`
	Type  string `json:"type"`
}

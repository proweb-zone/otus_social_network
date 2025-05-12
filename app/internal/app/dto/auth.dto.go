package dto

type AuthRequestDto struct {
	Email    string `json:"email" validate:"required,string"`
	Password string `json:"password" validate:"required,string"`
}

type AuthResponseDto struct {
	Token string `json:"token" validate:"required,string"`
}

package dto

type UsersRequestDto struct {
	First_name string `json:"first_name" validate:"required,string"`
	Last_name  string `json:"last_name" validate:"required,string"`
	Email      string `json:"email" validate:"required,string"`
	Password   string `json:"password" validate:"required,string"`
	Birth_date string `json:"birth_date" validate:"omitempty,string"`
	Gender     string `json:"gender" validate:"omitempty,string"`
	Hobby      string `json:"hobby" validate:"omitempty,string"`
	City       string `json:"city" validate:"omitempty,string"`
}

type UsersResponseDto struct {
	User_id uint `json:"user_id"`
}

func (d *UsersRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}

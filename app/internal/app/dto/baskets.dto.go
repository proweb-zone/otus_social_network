package dto

import (
	"strings"
)

type BasketsRequestDto struct {
	ProductId uint `json:"product_id" validate:"omitempty,numeric"`
	UserId    uint `json:"user_id" validate:"required,numeric"`
	Quantity  uint `json:"quantity" validate:"required,numeric"`
}

type UpdateBasketsRequestDto struct {
	ProductId uint `json:"product_id" validate:"omitempty,numeric"`
	Quantity  uint `json:"quantity" validate:"required,numeric"`
}

type QueryParamsDto struct {
	UserId    uint     `validate:"omitempty,numeric"`
	ProductId uint     `validate:"omitempty,numeric"`
	Sort      string   `validate:"omitempty"`
	Order     string   `validate:"omitempty,oneofci=asc desc"`
	Attrs     []string `validate:"omitempty"`
}

func (d *BasketsRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}

func (d *UpdateBasketsRequestDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
	}

	return validate(d, tagMsg).Error
}

func (d *QueryParamsDto) Validate() error {
	tagMsg := map[string]string{
		"required": "обязательно для заполнения",
		"numeric":  "должно быть числом",
		"oneofci":  "может быть одним из: ASC или DESC",
	}

	return validate(d, tagMsg).Error
}

func (d *QueryParamsDto) ValidateAttrs(validFields map[string]bool) []string {
	var result []string

	if len(d.Attrs) > 0 {
		for _, attr := range d.Attrs {
			if validFields[strings.ToLower(strings.Trim(attr, " "))] {
				result = append(result, attr)
			}
		}
	}

	return result
}

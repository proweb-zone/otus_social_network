package entity

import "time"

type Users struct {
	Id        uint      `json:"id" gorm:"column:id"`
	ProductId uint      `json:"product_id" gorm:"column:product_id"`
	UserId    uint      `json:"user_id" gorm:"column:user_id;index"`
	Quantity  uint      `json:"quantity" gorm:"column:quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Users) TableName() string {
	return "users"
}

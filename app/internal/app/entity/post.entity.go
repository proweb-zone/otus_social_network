package entity

import "time"

type Posts struct {
	ID          uint      `json:"id"`
	User_id     uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Posts) TableName() string {
	return "posts"
}

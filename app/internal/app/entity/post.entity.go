package entity

import "time"

type Posts struct {
	ID        uint      `json:"id"`
	User_id   uint      `json:"user_id"`
	Text      string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Posts) TableName() string {
	return "posts"
}

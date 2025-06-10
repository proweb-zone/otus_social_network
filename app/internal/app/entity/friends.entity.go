package entity

import "time"

type Friends struct {
	ID        uint      `json:"id"`
	User_id   uint      `json:"user_id`
	Friend_id uint      `json:"friend_id`
	CreatedAt time.Time `json:"created_at"`
}

func (Friends) TableName() string {
	return "friends"
}

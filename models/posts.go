package models

import "time"

type Post struct {
	Id        int64     `json:"id"`
	UserId    int64     `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAd time.Time `json:"created_at"`
}

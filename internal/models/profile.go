package models

import "time"

type Profile struct {
	UserID    int       `json:"user_id" db:"user_id"`
	Username  string    `json:"username" db:"username"`
	AvatarURL []byte    `json:"avatar" db:"avatar"`
	UserBIO   string    `json:"user_bio" db:"user_bio"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

package models

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email,omitempty" db:"email"`
	PasswordHash string    `json:"password_hash,omitempty" db:"password_hash"`
	Created_at   time.Time `json:"created_at,omitempty" db:"created_at"`
	Updated_at   time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Review struct {
	LandmarkID   int               `json:"landmark_id"`
	LandmarkName string            `json:"landmark_name"`
	UserID       int               `json:"user_id"`
	UserName     string            `json:"user_name"`
	Rating       int               `json:"rating"`
	Review       string            `json:"review"`
	Images       map[string][]byte `json:"images"`
}
type ReviewByUser struct {
	Username string `json:"username" db:"username"`
	Avatar   []byte `json:"avatar" db:"avatar"`
	Review
}

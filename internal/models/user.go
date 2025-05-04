package models

import "time"

type User struct {
	ID           int64     `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Created_at   time.Time `json:"created_at" db:"created_at"`
	Updated_at   time.Time `json:"updated_at" db:"updated_at"`
}

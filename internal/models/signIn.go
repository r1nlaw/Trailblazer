package models

import "time"

type SignInRequest struct {
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type JWTRequest struct {
	UserID    int64     `json:"user_id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

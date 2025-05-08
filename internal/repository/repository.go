package repository

import (
	"context"

	"trailblazer/internal/models"
)

type User interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, userData models.User) error
}

type Landmark interface {
	SaveLandmarks(ctx context.Context, landmarks []*models.Landmark) error
}
type Repository struct {
	User
	Landmark
}

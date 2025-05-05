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
	SaveLandmark(ctx context.Context, landmark *models.Landmark) error
}
type Repository struct {
	User
	Landmark
}

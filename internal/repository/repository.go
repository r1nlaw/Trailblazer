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
	GetFacilities(bbox models.BBOX) ([]models.Landmark, error)
	GetLandmarks(page int) ([]models.Landmark, error)
	GetLandmarksByIDs(ids []any) ([]models.Landmark, error)
	Search(q string) ([]models.Landmark, error)
	UpdateImagePath(place string, path string) error
	GetLandmarksByName(name string) (models.Landmark, error)
}
type Repository struct {
	User
	Landmark
}

package repository

import (
	"context"
	"trailblazer/internal/models"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, userData models.User) error
}

type Repository struct {
	User
}

func NewRepository(ctx context.Context, db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserPostgres(ctx, db),
	}
}

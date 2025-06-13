package service

import (
	"context"

	"trailblazer/internal/models"
	"trailblazer/internal/repository"
)

type User struct {
	repo repository.User
}

func (s *User) GetProfile(c context.Context, userID int64) (*models.Profile, error) {
	return s.repo.GetProfile(c, userID)
}

func NewUserService(repository repository.User) *User {
	return &User{
		repo: repository,
	}
}
func (s *User) GetUser(ctx context.Context, email string) (*models.User, error) {
	return s.repo.GetUser(ctx, email)
}

func (s *User) AddUser(c context.Context, user models.User) error {
	return s.repo.AddUser(c, user)
}

func (s *User) AddReview(review models.Review) error {
	return s.repo.AddReview(review)
}
func (s *User) GetReview(name string, onlyPhoto bool) (map[int]models.ReviewByUser, error) {
	return s.repo.GetReview(name, onlyPhoto)
}
func (s *User) UpdateUserProfile(c context.Context, i int, username string, bytes []byte, bio string) error {
	return s.repo.UpdateUserProfile(c, i, username, bytes, bio)
}

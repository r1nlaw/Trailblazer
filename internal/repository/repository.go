package repository

import (
	"context"

	"trailblazer/internal/models"
)

type User interface {
	GetUser(ctx context.Context, email string) (*models.User, error)
	AddUser(ctx context.Context, userData models.User) error
	UpdateUserProfile(context.Context, int, string, []byte, string) error
	GetProfile(ctx context.Context, userID int64) (*models.Profile, error)
}

type Landmark interface {
	GetFacilities(bbox models.BBOX) ([]models.Landmark, error)
	GetLandmarks(page int, categories []string) ([]models.Landmark, error)
	GetLandmarksByIDs(ids []any) ([]models.Landmark, error)
	Search(q string) ([]models.Landmark, error)
	UpdateImagePath(place string, path string) error
	GetLandmarksByName(name string) (models.Landmark, error)
	GetLandmarksByCategories(categories []string) ([]models.Landmark, error)
}
type Weather interface {
	SetWeather(id int, forecast models.WeatherForecast) error
	GetWeatherByLandmarkID(id int) (*[]models.WeatherResponse, error)
}
type Repository struct {
	User
	Weather
	Landmark
}

package service

import (
	"context"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
	"trailblazer/internal/utils"
)

type Service struct {
	repository *repository.Repository
	ctx        context.Context

	LandmarkService
	WeatherService
	UserService
}

type UserService interface {
	AddUser(c context.Context, user models.User) error
	GetUser(c context.Context, email string) (*models.User, error)
	AddReview(review models.Review) error
	GetReview(name string, onlyPhoto bool) (map[int]models.ReviewByUser, error)
	GetProfile(c context.Context, userID int64) (*models.Profile, error)
	UpdateUserProfile(c context.Context, i int, username string, bytes []byte, bio string) error
}
type LandmarkService interface {
	GetFacilities(bbox models.BBOX) ([]models.Landmark, error)
	GetLandmarks(page int, categories []string) ([]models.Landmark, error)
	GetLandmarksByIDs(ids []int) ([]models.Landmark, error)
	Search(q string) ([]models.Landmark, error)
	UpdateImagePath(place, path string) error
	GetLandmarksByName(name string) (models.Landmark, error)
	GetLandmarksByCategories(categories []string) ([]models.Landmark, error)
}
type WeatherService interface {
	SetWeather(id int, forecast models.WeatherForecast) error
	GetWeatherByLandmarkID(id int) (*[]models.WeatherResponse, error)
}

func NewService(ctx context.Context, repository *repository.Repository, tokenMaker utils.Maker, hashUtil utils.Hasher, cfg config.Config) *Service {
	return &Service{
		repository: repository,
		ctx:        ctx,

		LandmarkService: NewLandmarkService(repository.Landmark, cfg.ParserConfig),
		WeatherService:  NewWeatherService(repository.Weather, cfg.WeatherConfig),
		UserService:     NewUserService(repository.User),
	}
}

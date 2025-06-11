package service

import (
	"context"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
	"trailblazer/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repository *repository.Repository
	ctx        context.Context

	LandmarkService
	WeatherService
	UserService
}

type UserService interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	ChangeProfile(c *fiber.Ctx) error
	JWTMiddleware(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
	AddReview(review models.Review) error
	GetReview(name string, onlyPhoto bool) (map[int]models.ReviewByUser, error)
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
		UserService:     NewUserService(repository.User, tokenMaker, hashUtil),
	}
}

package service

import (
	"context"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
	"trailblazer/internal/service/hash"
	"trailblazer/internal/service/token"

	"github.com/gofiber/fiber/v2"
)

type Service struct {
	repository *repository.Repository
	ctx        context.Context
	tokenMaker token.Maker
	hashUtil   hash.Hasher
	LandmarkService
}

type UserService interface {
	SignUp(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
}
type LandmarkService interface {
	GetFacilities(bbox models.BBOX) ([]models.Landmark, error)
	GetLandmarks(page int) ([]models.Landmark, error)
}

func NewService(ctx context.Context, repository *repository.Repository, tokenMaker token.Maker, hashUtil hash.Hasher, cfg config.Config) *Service {
	return &Service{
		repository:      repository,
		ctx:             ctx,
		tokenMaker:      tokenMaker,
		hashUtil:        hashUtil,
		LandmarkService: NewLandmarkService(repository.Landmark, cfg.ParserConfig),
	}
}

func (s *Service) SignUp(c *fiber.Ctx) error {
	var request models.SignUpRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("failed to parse request")
	}

	hashedPassword, err := s.hashUtil.HashPassword(request.PasswordHash)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to hash password")
	}

	user := models.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.repository.AddUser(s.ctx, user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to add user")
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]string{"message": "user created successfully"})
}

func (s *Service) SignIn(c *fiber.Ctx) error {
	var request models.SignInRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("failed to parse request")
	}

	userData, err := s.repository.GetUser(s.ctx, request.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid email or password")
	}

	if !s.hashUtil.CheckPassword(userData.PasswordHash, request.PasswordHash) {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid email or password")
	}

	token, err := s.tokenMaker.CreateToken(userData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to create token")
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "login successful", "token": token})
}

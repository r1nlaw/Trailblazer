package service

import (
	"context"
	"encoding/base64"
	"log"
	"regexp"
	"strings"
	"time"

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
	ChangeProfile(c *fiber.Ctx) error
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
	if len(request.Username) > 30 || len(request.Username) < 2 {
		return c.Status(401).SendString("invalid username")
	}
	usernamePattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ]+(?: [a-zA-Zа-яА-ЯёЁ]+)*$`)
	if !usernamePattern.MatchString(request.Username) {
		return c.Status(402).SendString("username must contain only letters (latin or cyrillic)")
	}

	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.(com|ru|org|net|edu|gov|info|biz|io|me|dev)$`)

	if !emailRegex.MatchString(request.Email) {
		return c.Status(fiber.StatusBadRequest).SendString("invalid email address format")
	}
	if len(request.PasswordHash) < 6 {
		return c.Status(403).SendString("invalid password")
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
		return c.Status(405).SendString("invalid email or password")
	}

	if !s.hashUtil.CheckPassword(userData.PasswordHash, request.PasswordHash) {
		return c.Status(405).SendString("invalid email or password")
	}

	token, err := s.tokenMaker.CreateToken(userData.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to create token")
	}
	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false, // true, если HTTPS
		SameSite: "Lax", // или "None" при фронте на другом домене и HTTPS
		Path:     "/",
	})

	resp := SignInResponse{
		Message: "login successful",
		Token:   token,
		User:    *userData,
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

type SignInResponse struct {
	Message string      `json:"message"`
	Token   string      `json:"token"`
	User    models.User `json:"user"`
}

func (s *Service) ChangeProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization header format")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := s.tokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired token")
	}

	type incomingProfile struct {
		Username string `json:"username"`
		UserBIO  string `json:"user_bio"`
		Avatar   string `json:"avatar"`
	}

	var req incomingProfile
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("failed to parse request")
	}

	var avatarBytes []byte
	if req.Avatar != "" {
		avatarBytes, err = base64.StdEncoding.DecodeString(req.Avatar)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid base64 avatar")
		}
	}

	log.Printf("Parsed profile: username=%s, bio=%s, avatar len=%d, userID=%d",
		req.Username, req.UserBIO, len(avatarBytes), payload.UserID)

	err = s.repository.UpdateUserProfile(c.Context(), int(payload.UserID), req.Username, avatarBytes, req.UserBIO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to update profile: " + err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("profile updated successfully")
}

func (s *Service) GetUserProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization header format")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := s.tokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired token")
	}

	profile, err := s.repository.GetProfile(c.Context(), payload.UserID)
	if err != nil {
		log.Printf("Ошибка получения профиля пользователя (ID %d): %v", payload.UserID, err)
		return c.Status(fiber.StatusInternalServerError).SendString("failed to get user profile: " + err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(profile)
}

package service

import (
	"encoding/base64"
	"log"
	"regexp"
	"strings"
	"time"

	"trailblazer/internal/models"
	"trailblazer/internal/repository"
	"trailblazer/internal/utils"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	repo       repository.User
	TokenMaker utils.Maker
	hashUtil   utils.Hasher
}

func NewUserService(repository repository.User, tokenmaker utils.Maker, hashUtil utils.Hasher) *User {
	return &User{
		repo:       repository,
		TokenMaker: tokenmaker,
		hashUtil:   hashUtil,
	}
}

func (s *User) SignUp(c *fiber.Ctx) error {
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

	if err := s.repo.AddUser(c.Context(), user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to add user")
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]string{"message": "user created successfully"})
}

func (s *User) SignIn(c *fiber.Ctx) error {
	var request models.SignInRequest
	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("failed to parse request")
	}

	userData, err := s.repo.GetUser(c.Context(), request.Email)
	if err != nil {
		return c.Status(405).SendString("invalid email or password")
	}

	if !s.hashUtil.CheckPassword(userData.PasswordHash, request.PasswordHash) {
		return c.Status(405).SendString("invalid email or password")
	}

	token, err := s.TokenMaker.CreateToken(userData.ID)
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

func (s *User) ChangeProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization header format")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := s.TokenMaker.VerifyToken(tokenStr)
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

	err = s.repo.UpdateUserProfile(c.Context(), int(payload.UserID), req.Username, avatarBytes, req.UserBIO)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("failed to update profile: " + err.Error())
	}

	return c.Status(fiber.StatusOK).SendString("profile updated successfully")
}

func (s *User) GetUserProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization header format")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	payload, err := s.TokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired token")
	}

	profile, err := s.repo.GetProfile(c.Context(), payload.UserID)
	if err != nil {
		log.Printf("Ошибка получения профиля пользователя (ID %d): %v", payload.UserID, err)
		return c.Status(fiber.StatusInternalServerError).SendString("failed to get user profile: " + err.Error())
	}
	return c.Status(fiber.StatusOK).JSON(profile)
}

func (u *User) JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).SendString("missing authorization header")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid authorization header format")
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	_, err := u.TokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired token")
	}
	return c.Next()
}

func (s *User) AddReview(review models.Review) error {
	return s.repo.AddReview(review)
}
func (s *User) GetReview(name string, count int) (map[int]models.ReviewByUser, error) {
	return s.repo.GetReview(name, count)
}

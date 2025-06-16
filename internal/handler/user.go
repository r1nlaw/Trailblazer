package handler

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"trailblazer/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) SignUp(c *fiber.Ctx) error {
	var request models.SignUpRequest
	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "failed to parse request", err)
	}

	if len(request.Username) > 30 || len(request.Username) < 2 {
		return sendError(c, fiber.StatusBadRequest, "invalid username length", fmt.Errorf("username must be between 2 and 30 characters"))
	}

	usernamePattern := regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ]+(?: [a-zA-Zа-яА-ЯёЁ]+)*$`)
	if !usernamePattern.MatchString(request.Username) {
		return sendError(c, fiber.StatusBadRequest, "invalid username format", fmt.Errorf("username must contain only letters (latin or cyrillic)"))
	}

	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.(com|ru|org|net|edu|gov|info|biz|io|me|dev)$`)
	if !emailRegex.MatchString(request.Email) {
		return sendError(c, fiber.StatusBadRequest, "invalid email format", fmt.Errorf("invalid email address format"))
	}

	if len(request.PasswordHash) < 6 {
		return sendError(c, fiber.StatusBadRequest, "invalid password", fmt.Errorf("password must be at least 6 characters long"))
	}

	hashedPassword, err := h.hashUtil.HashPassword(request.PasswordHash)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to hash password", err)
	}

	user := models.User{
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
	}

	if err := h.service.UserService.AddUser(c.Context(), user); err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to add user", err)
	}

	if err := h.service.SendToken(request.Email); err != nil {
		if delErr := h.service.UserService.Delete(request.Email); delErr != nil {
			return sendError(c, fiber.StatusInternalServerError, "failed to delete user after token error", delErr)
		}
		return sendError(c, fiber.StatusInternalServerError, "failed to send verification token", err)
	}

	return c.Status(fiber.StatusCreated).JSON(map[string]string{"message": "user created successfully"})
}

func (h *Handler) SignIn(c *fiber.Ctx) error {
	var request models.SignInRequest
	if err := c.BodyParser(&request); err != nil {
		return sendError(c, fiber.StatusBadRequest, "failed to parse request", err)
	}

	userData, err := h.service.UserService.GetUser(c.Context(), request.Email)
	if err != nil {
		if err.Error() == "user is not verified" {
			return sendError(c, fiber.StatusUnauthorized, "user is not verified", err)
		}
		return sendError(c, fiber.StatusUnauthorized, "invalid credentials", fmt.Errorf("invalid email or password"))
	}

	if !h.hashUtil.CheckPassword(userData.PasswordHash, request.PasswordHash) {
		return sendError(c, fiber.StatusUnauthorized, "invalid credentials", fmt.Errorf("invalid email or password"))
	}

	token, err := h.TokenMaker.CreateToken(userData.ID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to create token", err)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    token,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
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

func (h *Handler) ChangeProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return sendError(c, fiber.StatusUnauthorized, "missing authorization header", fmt.Errorf("authorization header is required"))
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return sendError(c, fiber.StatusUnauthorized, "invalid authorization header format", fmt.Errorf("authorization header must start with 'Bearer '"))
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	payload, err := h.TokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return sendError(c, fiber.StatusUnauthorized, "invalid or expired token", err)
	}

	type incomingProfile struct {
		Username string `json:"username"`
		UserBIO  string `json:"user_bio"`
		Avatar   string `json:"avatar"`
	}

	var req incomingProfile
	if err := c.BodyParser(&req); err != nil {
		return sendError(c, fiber.StatusBadRequest, "failed to parse request", err)
	}

	var avatarBytes []byte
	if req.Avatar != "" {
		avatarBytes, err = base64.StdEncoding.DecodeString(req.Avatar)
		if err != nil {
			return sendError(c, fiber.StatusBadRequest, "invalid base64 avatar", err)
		}
	}

	log.Printf("Parsed profile: username=%s, bio=%s, avatar len=%d, userID=%d",
		req.Username, req.UserBIO, len(avatarBytes), payload.UserID)

	err = h.service.UserService.UpdateUserProfile(c.Context(), int(payload.UserID), req.Username, avatarBytes, req.UserBIO)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to update profile", err)
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "profile updated successfully"})
}

func (h *Handler) GetUserProfile(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return sendError(c, fiber.StatusUnauthorized, "missing authorization header", fmt.Errorf("authorization header is required"))
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return sendError(c, fiber.StatusUnauthorized, "invalid authorization header format", fmt.Errorf("authorization header must start with 'Bearer '"))
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	payload, err := h.TokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return sendError(c, fiber.StatusUnauthorized, "invalid or expired token", err)
	}

	profile, err := h.service.UserService.GetProfile(c.Context(), payload.UserID)
	if err != nil {
		return sendError(c, fiber.StatusInternalServerError, "failed to get user profile", err)
	}

	return c.Status(fiber.StatusOK).JSON(profile)
}

func (h *Handler) JWTMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return sendError(c, fiber.StatusUnauthorized, "missing authorization header", fmt.Errorf("authorization header is required"))
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return sendError(c, fiber.StatusUnauthorized, "invalid authorization header format", fmt.Errorf("authorization header must start with 'Bearer '"))
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := h.TokenMaker.VerifyToken(tokenStr)
	if err != nil {
		return sendError(c, fiber.StatusUnauthorized, "invalid or expired token", err)
	}

	return c.Next()
}

func (h *Handler) Verify(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return sendError(c, fiber.StatusBadRequest, "missing token", fmt.Errorf("verification token is required"))
	}

	err := h.service.VerifyEmail(token)
	if err != nil {
		return sendError(c, fiber.StatusUnauthorized, "invalid or expired token", err)
	}

	return c.Status(fiber.StatusOK).JSON(map[string]string{"message": "email verified successfully"})
}

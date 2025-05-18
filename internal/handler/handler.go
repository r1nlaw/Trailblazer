package handler

import (
	"log/slog"
	"os"

	"trailblazer/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
		AllowHeaders: "Content-Type, Authorization",
	}))
	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		slog.Info("Директория ./images не существует")
	}

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} - ${ua}\\n\n",
	}))
	app.Static("/images", "images")
	user := app.Group("/user")
	user.Post("/signIn", h.service.SignIn)
	user.Post("/signUp", h.service.SignUp)
	api := app.Group("/api")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	})).Post("/facilities", h.facilities)
	api.Get("/landmark", h.getLandmarks)
	api.Post("/getLandmarks", h.getLandmarksByIDs)
	api.Get("/search", h.search)
	api.Get("/landmark/:name", h.getLandmarksByName)
}

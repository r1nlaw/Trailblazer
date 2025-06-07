package handler

import (
	"log/slog"
	"os"

	"trailblazer/internal/api"
	"trailblazer/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	service *service.Service
	api     api.WeatherAPI
}

func NewHandler(service *service.Service, api api.WeatherAPI) *Handler {
	return &Handler{service: service, api: api}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))

	if _, err := os.Stat("./images"); os.IsNotExist(err) {
		slog.Info("Директория ./images не существует")
	}

	app.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} - ${ua}\\n\n",
	}))
	app.Static("/assets", "assets")
	app.Static("/images", "images")
	user := app.Group("/user")
	user.Post("/signIn", h.service.SignIn)
	user.Post("/signUp", h.service.SignUp)
	user.Post("/changeProfile", h.service.ChangeProfile)
	user.Get("/profile", h.service.GetUserProfile)
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

package handler

import (
	"trailblazer/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173/",
		AllowMethods:     "GET, POST, PUT, DELETE",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true,
	}))
	user := app.Group("/user")
	user.Post("/signIn", h.service.SignIn)
	user.Post("/signUp", h.service.SignUp)
	api := app.Group("/api")
	api.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	})).Post("/facilities", h.facilities)
	api.Get("/landmark", h.getLandmarks)
}

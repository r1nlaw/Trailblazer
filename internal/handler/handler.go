package handler

import (
	"log/slog"
	"os"

	"trailblazer/internal/api"
	"trailblazer/internal/service"
	"trailblazer/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type Handler struct {
	service    *service.Service
	api        api.WeatherAPI
	TokenMaker utils.Maker
	hashUtil   utils.Hasher
}

func NewHandler(service *service.Service, api api.WeatherAPI, hashutil utils.Hasher, maker utils.Maker) *Handler {
	return &Handler{service: service, api: api, TokenMaker: maker, hashUtil: hashutil}
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
	app.Static("/resources", "resources")
	app.Static("/images", "images")
	user := app.Group("/user")
	user.Post("/signIn", h.SignIn)
	user.Post("/signUp", h.SignUp)
	user.Post("/changeProfile", h.ChangeProfile)
	user.Get("/profile", h.GetUserProfile)
	review := user.Group("/review")
	review.Use("/add/:name", h.JWTMiddleware)
	review.Post("/add/:name", h.AddReview)
	review.Get("/get/:name", h.GetReview)
	apiGroup := app.Group("/api")
	apiGroup.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: false,
	})).Post("/facilities", h.facilities)
	apiGroup.Get("/landmark", h.getLandmarks)
	apiGroup.Post("/getLandmarks", h.getLandmarksByIDs)
	apiGroup.Get("/search", h.search)
	apiGroup.Get("/landmark/:name", h.getLandmarksByName)

}

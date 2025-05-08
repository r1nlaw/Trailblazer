package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trailblazer/internal/config"
	"trailblazer/internal/handler"
	"trailblazer/internal/repository"
	"trailblazer/internal/service"
	"trailblazer/internal/service/hash"
	"trailblazer/internal/service/token"

	"github.com/gofiber/fiber/v2"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
func main() {
	logger := InitLogger()
	slog.SetDefault(logger)

	configPath := flag.String("c", "configs/config.yml", "The path to the configuration file")
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("error to parse config: %v", err))
	}
	var repo *repository.Repository
	repo, err = repository.NewPostgresRepository(context.Background(), cfg.DatabaseConfig)
	if err != nil {
		slog.Warn("failed to initialize DB", err)
		return
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		slog.Warn("JWT_SECRET must be set", err)
		return
	}
	tokenMaker, err := token.NewJWTMaker(secret)
	if err != nil {
		slog.Warn("failed to initialize tokenMaker", err)
		return
	}

	hashUtil := hash.NewBcryptHasher()
	ctx := context.Background()
	slog.Info("initializing repository")
	services := service.NewService(ctx, repo, tokenMaker, hashUtil)
	slog.Info("initializing services")
	handlers := handler.NewHandler(services)
	app := fiber.New()

	handlers.InitRoutes(app)
	go func() {
		if err := app.Listen(":" + cfg.HostConfig.Port); err != nil {
			slog.Warn("error starting server", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		slog.Warn("error shutting down server", err)
	}

	slog.Info("Server stopped successfully")
}

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

	api2 "trailblazer/internal/api"
	"trailblazer/internal/config"
	"trailblazer/internal/handler"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
	"trailblazer/internal/service"
	"trailblazer/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/sabloger/sitemap-generator/smg"
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
	tokenMaker, err := utils.NewJWTMaker(secret)
	if err != nil {
		slog.Warn("failed to initialize tokenMaker", err)
		return
	}

	hashUtil := utils.NewBcryptHasher()
	ctx := context.Background()
	api := api2.NewWeatherClient(cfg.WeatherConfig)
	slog.Info("initializing repository")
	services := service.NewService(ctx, repo, tokenMaker, hashUtil, *cfg)
	slog.Info("initializing services")
	handlers := handler.NewHandler(services, *api, hashUtil, tokenMaker)
	app := fiber.New()

	handlers.InitRoutes(app)
	go func() {
		if err := app.Listen(":" + cfg.HostConfig.Port); err != nil {
			slog.Warn("error starting server", err)
		}
	}()
	go func() {
		ticker := time.NewTicker(5 * 24 * time.Hour)
		landmarks, err := repo.Landmark.GetLandmarks(-1, nil)
		if err != nil {
			slog.Warn("failed to get landmarks: ", err)
			return
		}
		CreateSiteMap(landmarks, "resources", "https://putevod-crimea.ru")

		for _ = range ticker.C {
			landmarks, err := repo.Landmark.GetLandmarks(-1, nil)
			if err != nil {
				slog.Warn("failed to get landmarks: ", err)
				return
			}
			CreateSiteMap(landmarks, "assets", "https://putevod-crimea.ru")

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

func CreateSiteMap(landmarks []models.Landmark, path, domain string) {
	now := time.Now().UTC()
	sm := smg.NewSitemap(true)
	sm.SetName("sitemap")
	sm.SetHostname(domain)
	sm.SetOutputPath(path)
	sm.SetLastMod(&now)
	sm.SetCompress(false)
	sm.SetMaxURLsCount(50000)

	for _, url := range landmarks {
		err := sm.Add(&smg.SitemapLoc{
			Loc:        fmt.Sprintf("%s/landmark/%s", domain, url.TranslatedName),
			LastMod:    &now,
			ChangeFreq: smg.Always,
			Priority:   0.7,
		})
		if err != nil {
			slog.Warn(fmt.Sprintf("add sitemap loc err: %v", err))
		}
	}
	err := sm.Add(&smg.SitemapLoc{
		Loc:        fmt.Sprintf("%s/", domain),
		LastMod:    &now,
		ChangeFreq: smg.Always,
		Priority:   1,
	})
	if err != nil {
		slog.Warn(fmt.Sprintf("add sitemap loc err: %v", err))
	}
	filenames, err := sm.Save()
	if err != nil {
		slog.Error("Unable to Save Sitemap:", err)
		os.Exit(1)
	}
	for i, filename := range filenames {
		slog.Info(fmt.Sprintf("file no.%d %s", i+1, filename))
	}
}

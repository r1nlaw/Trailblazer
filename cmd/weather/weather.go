package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	api2 "trailblazer/internal/api"
	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
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
	api := api2.NewWeatherClient(cfg.WeatherConfig)

	if err != nil {
		slog.Error(fmt.Sprintf("error to get location: %v", err))
	}
	//fmt.Println(res)
	var repo *repository.Repository
	repo, err = repository.NewPostgresRepository(context.Background(), cfg.DatabaseConfig)
	if err != nil {
		slog.Warn("failed to initialize DB", err)
		return
	}
	slog.Info("initializing repository")
	landmarks, err := repo.Landmark.GetLandmarks(-1, nil)
	if err != nil {
		slog.Warn("failed to get landmarks: ", err)
		return
	}

	ticker := time.NewTicker(8 * time.Hour).C

	for _, landmark := range landmarks {
		slog.Info(fmt.Sprintf("%d-%s", landmark.ID, landmark.Name))
		res, err := api.WeatherAt(models.Location{
			Lng: landmark.Lng,
			Lat: landmark.Lat,
		})
		if err != nil {
			slog.Warn("failed to get weather at: ", landmark, err)
			continue
		}
		repo.Weather.SetWeather(landmark.ID, res)

	}
	slog.Info("Цикл завершен")
	for _ = range ticker {
		for _, landmark := range landmarks {
			slog.Info(fmt.Sprintf("%d-%s", landmark.ID, landmark.Name))

			res, err := api.WeatherAt(models.Location{
				Lng: landmark.Lng,
				Lat: landmark.Lat,
			})
			if err != nil {
				slog.Warn("failed to get weather at: ", landmark, err)
				continue
			}
			repo.Weather.SetWeather(landmark.ID, res)

		}
		slog.Info("Цикл завершен")
	}

}

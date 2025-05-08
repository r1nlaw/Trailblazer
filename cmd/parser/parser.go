package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"trailblazer/internal/config"
	"trailblazer/internal/repository"
	"trailblazer/internal/service"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
func main() {
	logger := InitLogger()
	slog.SetDefault(logger)

	configPath := flag.String("c", "./configs/config.yml", "The path to the configuration file")
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("error to parse config: %v", err))
	}
	var repo *repository.Repository
	repo, err = repository.NewJSONRepository(cfg.Dir)
	if err != nil {
		slog.Warn("failed to initialize DB", err)
		return
	}
	var services *service.Service
	services = service.NewService(context.Background(), repo, nil, nil, *cfg)
	services.Crawl()
	fmt.Println(repo)

}

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"trailblazer/internal/config"
	"trailblazer/internal/repository"
)

func InitLogger() *slog.Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	return logger
}
func main() {
	logger := InitLogger()
	slog.SetDefault(logger)

	configPath := flag.String("c", "config/config.yaml", "The path to the configuration file")
	flag.Parse()
	cfg, err := config.New(*configPath)
	if err != nil {
		slog.Error(fmt.Sprintf("error to parse config: %v", err))
	}
	var repo *repository.Repository

	repo, err = repository.NewJSONRepository(".")
	if err != nil {
		slog.Warn("failed to initialize DB", err)
		return
	}

}

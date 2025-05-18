package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"trailblazer/internal/config"
	"trailblazer/internal/repository"
	"trailblazer/internal/service"
	"trailblazer/internal/service/hash"
	"trailblazer/internal/service/token"

	"github.com/mozillazg/go-unidecode"
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
	services := service.NewService(ctx, repo, tokenMaker, hashUtil, *cfg)
	slog.Info("initializing services")

	dirFiles, err := os.ReadDir("images")
	if err != nil {
		panic(err)
	}
	for _, file := range dirFiles {
		fileName := file.Name()
		ext := filepath.Ext(fileName)
		fileName = strings.Replace(fileName, ".jpg", "", -1)
		fileName = strings.Replace(fileName, ".jpeg", "", -1)
		translator := unidecode.Unidecode(fileName)
		replacer := strings.NewReplacer(
			`-`, `_`,
			`.`, ``,
			` `, `_`,
			`,`, ``,
			`«`, ``,
			`»`, ``,
			`'`, ``,
			`>`, ``,
			`<`, ``,
			`"`, ``,
		)
		translator = replacer.Replace(translator)

		err = services.LandmarkService.UpdateImagePath(fileName, translator+ext)

		if err != nil {
			slog.Error("failed to update image path", err)
		}
		err = os.Rename("images/"+fileName+ext, "images/"+translator+ext)
		if err != nil {
			slog.Error("failed to rename image file", err)
		}
		fmt.Println(translator)

	}
}

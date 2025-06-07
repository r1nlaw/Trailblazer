package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/sabloger/sitemap-generator/smg"

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
	CreateSiteMap(landmarks, "assets", "https://putevod-crimea.ru")

}

type URL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
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

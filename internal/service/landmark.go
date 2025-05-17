package service

import (
	"net/http"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
)

type Landmark struct {
	repo repository.Landmark
	config.ParserConfig
	cookies   []*http.Cookie
	ApiConfig config.GeocoderConfig
}

func NewLandmarkService(landmark repository.Landmark, cfg config.ParserConfig) *Landmark {
	return &Landmark{
		repo:         landmark,
		ParserConfig: cfg,
	}
}

func (s *Landmark) GetFacilities(bbox models.BBOX) ([]models.Landmark, error) {
	return s.repo.GetFacilities(bbox)

}

func (s *Landmark) GetLandmarks(page int) ([]models.Landmark, error) {
	return s.repo.GetLandmarks(page)
}

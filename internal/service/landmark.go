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
	cookies []*http.Cookie
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

func (s *Landmark) GetLandmarks(page int, categories []string) ([]models.Landmark, error) {
	return s.repo.GetLandmarks(page, categories)
}

func (s *Landmark) GetLandmarksByIDs(ids []int) ([]models.Landmark, error) {
	ID := make([]any, len(ids))
	for i, id := range ids {
		ID[i] = id
	}
	return s.repo.GetLandmarksByIDs(ID)

}
func (s *Landmark) Search(q string) ([]models.Landmark, error) {
	return s.repo.Search(q)
}

func (s *Landmark) UpdateImagePath(place, path string) error {
	return s.repo.UpdateImagePath(place, path)
}
func (s *Landmark) GetLandmarksByName(name string) (models.Landmark, error) {
	return s.repo.GetLandmarksByName(name)
}

func (s *Landmark) GetLandmarksByCategories(categories []string) ([]models.Landmark, error) {
	return s.repo.GetLandmarksByCategories(categories)
}

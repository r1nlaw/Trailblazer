package repository

import (
	"context"
	"encoding/json"
	"os"

	"trailblazer/internal/models"
)

type LandmarkJSON struct {
	db *JSONRepository
}

func NewLandmarkJSON(db *JSONRepository) *LandmarkJSON {
	return &LandmarkJSON{
		db: db,
	}
}

func (r *LandmarkJSON) SaveLandmark(ctx context.Context, landmark models.Landmark) error {
	landmarks, err := r.LoadLandmarks(ctx)
	if err != nil {
		return err
	}

	landmarks = append(landmarks, landmark)

	data, err := json.MarshalIndent(landmarks, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(r.db.path, data, 0644); err != nil {
		return err
	}

	return nil
}

func (r *LandmarkJSON) LoadLandmarks(ctx context.Context) ([]models.Landmark, error) {
	file, err := os.ReadFile(r.db.path)
	if err != nil {
		return nil, err
	}
	var landmarks []models.Landmark
	if err := json.Unmarshal(file, &landmarks); err != nil {
		return nil, err
	}
	return landmarks, nil
}

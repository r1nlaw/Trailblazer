package repository

import (
	"log/slog"
	"os"
)

type JSONRepository struct {
	path string
}

func newJsonDb(path string) (*JSONRepository, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("[]"), 0644); err != nil {
			slog.Warn("failed to initialize json db: " + err.Error())
			return nil, err
		}
	} else if err != nil {
		slog.Warn("failed to stat file: " + err.Error())
		return nil, err
	}
	return &JSONRepository{path: path}, nil
}

func NewJSONRepository(path string) (*Repository, error) {
	db, err := newJsonDb(path)
	if err != nil {
		return nil, err
	}
	landmarkDb := NewLandmarkJSON(db)
	return &Repository{
		Landmark: landmarkDb,
	}, nil
}

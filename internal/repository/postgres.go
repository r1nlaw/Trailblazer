package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"trailblazer/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {

	if err := godotenv.Load("./.env"); err != nil {
		log.Fatal("error to read env file")
	}

	connectStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
	)

	db, err := sqlx.Connect("postgres", connectStr)

	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://./schema",
		cfg.DBName, driver,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return db, nil
}

func NewPostgresRepository(ctx context.Context, cfg config.DatabaseConfig) (*Repository, error) {
	db, err := NewPostgresDB(cfg)
	if err != nil {
		return nil, err
	}
	userDb := NewUserPostgres(ctx, db)
	landmarkDB := NewLandmarkPostgres(ctx, db)
	repository := &Repository{
		User:     userDb,
		Landmark: landmarkDB,
	}
	return repository, nil
}

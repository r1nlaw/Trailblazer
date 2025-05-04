package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"trailblazer/internal/handler"
	"trailblazer/internal/repository"
	"trailblazer/internal/service"
	"trailblazer/internal/service/hash"
	"trailblazer/internal/service/token"
	logging "trailblazer/pkg"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("./configs/config.yml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		logging.Logger.Warn(logging.MakeLog("failed to load env file", err))
		return
	}

	if err := logging.NewLogService(os.Stdout, os.Getenv("LOG_MODE")); err != nil {
		log.Fatal("Failed to initialize logger: ", err)
	}

	db, err := repository.InitDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logging.Logger.Warn(logging.MakeLog("failed to initialize DB", err))
		return
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logging.Logger.Warn(logging.MakeLog("JWT_SECRET must be set", err))
		return
	}
	tokenMaker, err := token.NewJWTMaker(secret)
	if err != nil {
		logging.Logger.Warn(logging.MakeLog("failed to initialize tokenMaker", err))
		return
	}

	hashUtil := hash.NewBcryptHasher()
	ctx := context.Background()
	repo := repository.NewRepository(ctx, db)
	logging.Logger.Info("initializing repository")
	services := service.NewService(ctx, repo, tokenMaker, hashUtil)
	logging.Logger.Info("initializing services")
	handlers := handler.NewHandler(services)
	app := fiber.New()

	handlers.InitRoutes(app)
	go func() {
		if err := app.Listen(":" + viper.GetString("server.port")); err != nil {
			logging.Logger.Warn(logging.MakeLog("error starting server", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		logging.Logger.Warn(logging.MakeLog("error shutting down server", err))
	}

	logging.Logger.Info("Server stopped successfully")
}

package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	DatabaseConfig
	HostConfig
	GeocoderConfig
	ParserConfig
}
type HostConfig struct {
	Port string
}

type DatabaseConfig struct {
	Dir      string
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}
type GeocoderConfig struct {
	GeoUrl string
	ApiKey string
}
type ParserConfig struct {
	IsProduction bool
	BaseURL      string
}

func New(path string) (*Config, error) {
	viper.SetConfigFile(path)
	//res, _ := os.ReadDir(".")
	//for _, r := range res {
	//	fmt.Println(r.Name())
	//}
	if err := viper.ReadInConfig(); err != nil {
		slog.Warn(fmt.Sprintf("Error reading config: %v", err))
		return nil, err
	}
	var cfg Config
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load env file", err)
	}
	cfg.DatabaseConfig = DatabaseConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Dir:      viper.GetString("db.dir"),
	}
	cfg.HostConfig = HostConfig{
		Port: viper.GetString("server.port"),
	}
	cfg.GeocoderConfig = GeocoderConfig{
		GeoUrl: viper.GetString("geocoder.url"),
		ApiKey: viper.GetString("geocoder.api_key"),
	}
	cfg.ParserConfig = ParserConfig{
		IsProduction: viper.GetBool("parser.is_production"),
		BaseURL:      viper.GetString("parser.base_url"),
	}
	return &cfg, nil
}

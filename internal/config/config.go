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
	WeatherConfig
	ParserConfig
	SMTPConfig
}
type HostConfig struct {
	Port   string
	Domain string
}

type SMTPConfig struct {
	Email    string
	Password string
	Username string
	Host     string
	Port     string
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
type WeatherConfig struct {
	WeatherUrl string
	ApiKey     string
	Language   string
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
		Port:   viper.GetString("server.port"),
		Domain: viper.GetString("server.domain"),
	}
	cfg.WeatherConfig = WeatherConfig{
		WeatherUrl: viper.GetString("weather.url"),
		ApiKey:     os.Getenv("WEATHER_API"),
		Language:   viper.GetString("weather.lang"),
	}
	cfg.ParserConfig = ParserConfig{
		IsProduction: viper.GetBool("parser.is_production"),
		BaseURL:      viper.GetString("parser.base_url"),
	}
	cfg.SMTPConfig = SMTPConfig{
		Email:    os.Getenv("EMAIL_FROM"),
		Password: os.Getenv("SMTP_PASS"),
		Username: os.Getenv("SMTP_USER"),
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
	}
	return &cfg, nil
}

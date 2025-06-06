package service

import (
	"trailblazer/internal/config"
	"trailblazer/internal/models"
	"trailblazer/internal/repository"
)

type Weather struct {
	repo repository.Weather
	config.WeatherConfig
}

func (w Weather) SetWeather(id int, forecast models.WeatherForecast) error {
	return w.repo.SetWeather(id, forecast)
}

func (w Weather) GetWeatherByLandmarkID(id int) (*[]models.WeatherResponse, error) {
	return w.repo.GetWeatherByLandmarkID(id)
}

func NewWeatherService(weather repository.Weather, weatherConfig config.WeatherConfig) Weather {
	return Weather{
		repo:          weather,
		WeatherConfig: weatherConfig,
	}
}

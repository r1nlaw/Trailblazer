package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"trailblazer/internal/config"
	"trailblazer/internal/models"
)

type WeatherAPI struct {
	baseURL string
	token   string
	Client  *http.Client
	lang    string
}

func NewWeatherClient(cfg config.WeatherConfig) *WeatherAPI {
	return &WeatherAPI{
		baseURL: cfg.WeatherUrl,
		token:   cfg.ApiKey,
		Client:  http.DefaultClient,
		lang:    cfg.Language,
	}
}

func (c *WeatherAPI) WeatherAt(loc models.Location) (models.WeatherForecast, error) {
	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", loc.Lat))
	params.Set("lon", fmt.Sprintf("%f", loc.Lng))
	params.Set("units", "metric")
	params.Set("lang", c.lang)
	params.Set("appid", c.token)
	req, err := http.NewRequest(http.MethodGet, c.baseURL+"/data/2.5/forecast?"+params.Encode(), nil)
	if err != nil {
		return models.WeatherForecast{}, err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.Client.Do(req)
	if err != nil {
		return models.WeatherForecast{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return models.WeatherForecast{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	var weatherResponse models.WeatherForecast
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return models.WeatherForecast{}, err
	}
	if err := json.Unmarshal(body, &weatherResponse); err != nil {
		return models.WeatherForecast{}, err
	}
	return weatherResponse, nil
}

package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"trailblazer/internal/models"

	"github.com/jmoiron/sqlx"
)

type WeatherDB struct {
	ctx      context.Context
	postgres *sqlx.DB
}

func NewWeatherPostgres(ctx context.Context, db *sqlx.DB) *WeatherDB {
	return &WeatherDB{ctx: ctx, postgres: db}
}

func (r *WeatherDB) SetWeather(id int, forecast models.WeatherForecast) error {
	query := `	
			INSERT INTO weather(landmark_id, date, temperature, description, icon,rain,wind_speed,wind_degree )
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			ON CONFLICT (landmark_id,date) DO UPDATE SET temperature=$3, description=$4, icon=$5,rain=$6,wind_speed=$7,wind_degree=$8
				`
	for _, info := range forecast.List {
		date := time.Unix(info.Dt, 0)
		temperature := info.Main.Temp
		var description, icon string
		if len(info.Weather) > 0 {
			description = info.Weather[0].Description
			icon = info.Weather[0].Icon
		}
		var rain, windSpeed float64
		var windDegree int
		if info.Rain != nil {
			rain = info.Rain.ThreeHour
		}
		windSpeed = info.Wind.Speed
		windDegree = info.Wind.Deg

		_, err := r.postgres.Exec(query, id, date, temperature, description, icon, rain, windSpeed, windDegree)
		if err != nil {
			slog.Error(fmt.Sprintf("error with saving weather %v", err))
			continue
		}

	}
	return nil
}

func (r *WeatherDB) GetWeatherByLandmarkID(id int) (*[]models.WeatherResponse, error) {
	query := `
			SELECT date,temperature,description,icon,rain,wind_speed,wind_degree FROM weather WHERE 
				landmark_id=$1 AND date>current_timestamp 
 			`
	res, err := r.postgres.Query(query, id)
	if err != nil {
		return nil, err
	}
	weatherList := make([]models.WeatherResponse, 0)
	for res.Next() {
		weather := new(models.WeatherResponse)
		weather.LandmarkID = id
		err = res.Scan(&weather.Date, &weather.Temperature, &weather.Description, &weather.Icon, &weather.Rain, &weather.WindSpeed, &weather.WindDegree)
		if err != nil {
			return nil, err
		}
		weatherList = append(weatherList, *weather)
	}
	return &weatherList, nil

}

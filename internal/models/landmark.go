package models

import (
	"time"
)

type Landmark struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	TranslatedName  string     `json:"translated_name"`
	Address         string     `json:"address"`
	Category        string     `json:"category"`
	Schedules       []Schedule `json:"schedules"`
	Prices          []Price    `json:"prices"`
	Description     string     `json:"description"`
	History         string     `json:"history"`
	Location        `json:"location"`
	ImagePath       string             `json:"image_path"`
	WeatherResponse *[]WeatherResponse `json:"weathers"`
}
type Schedule struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
}
type Price struct {
	Value       float64 `json:"value"`
	Currency    string  `json:"currency"`
	Description string  `json:"description"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
type BBOX struct {
	SW Point `json:"sw"`
	NE Point `json:"ne"`
}

type Point struct {
	Lng float64 `json:"lng"`
	Lat float64 `json:"lat"`
}

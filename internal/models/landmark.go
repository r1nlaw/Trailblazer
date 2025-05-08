package models

import (
	"time"
)

type Landmark struct {
	Name        string     `json:"name"`
	Address     string     `json:"address"`
	Schedules   []Schedule `json:"schedules"`
	Prices      []Price    `json:"prices"`
	Description string     `json:"description"`
	History     string     `json:"history"`
	Location    `json:"location"`
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

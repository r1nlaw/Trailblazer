package utils

import (
	"strconv"
	"strings"

	"trailblazer/internal/models"
)

func LocationFromPoint(p string) models.Location {
	p = strings.TrimPrefix(p, "POINT(")
	p = strings.TrimSuffix(p, ")")
	numbers := strings.Fields(p)
	lon, _ := strconv.ParseFloat(numbers[0], 64)
	lat, _ := strconv.ParseFloat(numbers[1], 64)
	loc := models.Location{
		Lat: lat,
		Lng: lon,
	}
	return loc
}

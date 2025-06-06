package handler

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"trailblazer/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) facilities(c *fiber.Ctx) error {
	var req models.BBOX
	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest)
		_, _ = c.WriteString(err.Error())
	}
	facilities, err := h.service.GetFacilities(req)
	for i := range facilities {
		facilities[i].WeatherResponse, err = h.service.WeatherService.GetWeatherByLandmarkID(facilities[i].ID)
		if err != nil {
			slog.Error(fmt.Sprintf("error with finding weather %v", err))
			continue
		}
	}
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(facilities)
}
func (h *Handler) getLandmarks(ctx *fiber.Ctx) error {
	var page int
	p := ctx.Query("page", "1")
	var categories []string
	ctx.Request().URI().QueryArgs().VisitAll(func(key, val []byte) {
		if string(key) == "category" {
			categories = append(categories, fmt.Sprintf("'%s'", strings.ToLower(string(val))))
		}
	})
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
		err = nil
	}
	landmarks, err := h.service.LandmarkService.GetLandmarks(page, categories)
	for i := range landmarks {
		landmarks[i].WeatherResponse, err = h.service.WeatherService.GetWeatherByLandmarkID(landmarks[i].ID)
		if err != nil {
			slog.Error(fmt.Sprintf("error with finding weather %v", err))
			continue
		}
	}
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(landmarks)

}

func (h *Handler) getLandmarksByIDs(ctx *fiber.Ctx) error {
	var req []int
	if err := ctx.BodyParser(&req); err != nil {

		return ctx.JSON(fiber.Map{"error": err.Error()})
	}
	points, err := h.service.GetLandmarksByIDs(req)
	for i := range points {
		points[i].WeatherResponse, err = h.service.WeatherService.GetWeatherByLandmarkID(points[i].ID)
		if err != nil {
			slog.Error(fmt.Sprintf("error with finding weather %v", err))
			continue
		}
	}
	if err != nil {
		return ctx.JSON(fiber.Map{"error": err.Error})
	}
	return ctx.JSON(points)
}

func (h *Handler) search(ctx *fiber.Ctx) error {
	query := ctx.Query("q", "1")
	landmarks, err := h.service.LandmarkService.Search(query)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return ctx.JSON(landmarks)
}

func (h *Handler) getLandmarksByName(ctx *fiber.Ctx) error {
	name := ctx.Params("name")

	landmark, err := h.service.LandmarkService.GetLandmarksByName(name)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	landmark.WeatherResponse, err = h.service.WeatherService.GetWeatherByLandmarkID(landmark.ID)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err})
	}
	return ctx.JSON(landmark)
}

func (h *Handler) getLandmarksByCategories(ctx *fiber.Ctx) error {
	var categories []string
	ctx.Request().URI().QueryArgs().VisitAll(func(key, val []byte) {
		if string(key) == "category" {
			categories = append(categories, string(val))
		}
	})

	if len(categories) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "no categories provided"})
	}

	landmarks, err := h.service.LandmarkService.GetLandmarksByCategories(categories)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(landmarks)
}

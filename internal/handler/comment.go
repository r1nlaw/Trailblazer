package handler

import (
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"

	"trailblazer/internal/models"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) AddReview(c *fiber.Ctx) error {
	name := c.Params("name")
	landmarkIDStr := c.FormValue("landmark_id")
	landmarkID, err := strconv.Atoi(landmarkIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Недействительный landmark_id"})
	}
	userIDStr := c.FormValue("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Недействительный user_id"})
	}
	ratingStr := c.FormValue("rating")

	rating, err := strconv.Atoi(ratingStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Недействительная оценка"})
	}
	review := c.FormValue("review")
	images := make(map[string][]byte)
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Ошибка при получении формы"})
	}
	for key, files := range form.File {
		if strings.HasPrefix(key, "images[") && strings.HasSuffix(key, "]") {
			filename := strings.TrimPrefix(strings.TrimSuffix(key, "]"), "images[")
			if len(files) > 0 {
				file, err := files[0].Open()
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Ошибка при открытии файла %s", filename)})
				}
				defer file.Close()

				data, err := io.ReadAll(file)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Ошибка при чтении файла %s", filename)})
				}
				images[filename] = data
			}
		}
	}
	slog.Info(fmt.Sprintf(" %s: %s %d from %d", name, review, rating, userID))
	Review := models.Review{
		LandmarkID: landmarkID,
		UserID:     userID,
		Rating:     rating,
		Review:     review,
		Images:     images,
	}
	err = h.service.UserService.AddReview(Review)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"review": Review})
}

func (h *Handler) GetReview(c *fiber.Ctx) error {
	name := c.Params("name")
	countStr := c.Query("count")
	if countStr == "" {
		countStr = "0"
	}
	count, err := strconv.Atoi(countStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	reviews, err := h.service.UserService.GetReview(name, count)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"reviews": reviews})
}

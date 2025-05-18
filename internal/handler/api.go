package handler

import (
	"strconv"

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
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(facilities)
}
func (h *Handler) getLandmarks(ctx *fiber.Ctx) error {
	var page int
	p := ctx.Query("page", "1")
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
		err = nil
	}
	landmarks, err := h.service.LandmarkService.GetLandmarks(page)
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

package models

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Server struct {
	fiberApp *fiber.App
}

func (s *Server) Run(port string, handler fiber.Router) error {
	fmt.Printf("Server start on port: %s\n", port)

	s.fiberApp = fiber.New()

	s.fiberApp.Server().ReadTimeout = 10 * time.Second
	s.fiberApp.Server().WriteTimeout = 10 * time.Second

	return s.fiberApp.Listen(":" + port)
}

func (s *Server) Close(ctx context.Context) error {
	if err := s.fiberApp.Shutdown(); err != nil {
		return fmt.Errorf("server close error: %v", err)
	}
	return nil
}

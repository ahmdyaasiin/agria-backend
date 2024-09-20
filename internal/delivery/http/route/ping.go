package route

import (
	"github.com/gofiber/fiber/v3"
)

func NewPingRoute(route fiber.Router) {
	route.Get("/ping", func(c fiber.Ctx) error {
		return c.SendString("pong")
	})
	route.Get("/error", func(c fiber.Ctx) error {
		return fiber.NewError(fiber.StatusInternalServerError, "route for error test")
	})
}

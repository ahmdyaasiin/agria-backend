package route

import "github.com/gofiber/fiber/v3"

func NewAdminRoutes(route fiber.Router) {
	route.Get("/hello", func(ctx fiber.Ctx) error {
		return ctx.SendString("hello admin!")
	})
}

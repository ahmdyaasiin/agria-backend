package route

import "github.com/gofiber/fiber/v3"

func NewUserRoutes(route fiber.Router) {
	route.Get("/hello", func(ctx fiber.Ctx) error {
		return ctx.SendString("hello user!")
	})
}

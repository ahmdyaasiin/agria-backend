package http

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/route"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	App        *fiber.App
	Middleware fiber.Handler
	Cors       fiber.Handler
}

func (c *Config) Start() {
	v1 := c.App.Group("/v1")

	route.NewPingRoute(v1)
	route.NewUserRoutes(v1)
	route.NewAdminRoutes(v1.Group("/admin"))
}

package http

import (
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/handler/interfaces"
	"github.com/ahmdyaasiin/agria-backend/internal/delivery/http/route"
	"github.com/gofiber/fiber/v3"
)

type Config struct {
	App         *fiber.App
	Middleware  fiber.Handler
	Cors        fiber.Handler
	UserHandler interfaces.UserHandler
}

func (c *Config) Start() {
	v1 := c.App.Group("/v1")

	route.NewPingRoute(v1)
	route.NewUserRoutes(v1, c.UserHandler)
	route.NewAdminRoutes(v1.Group("/admin"))
}
